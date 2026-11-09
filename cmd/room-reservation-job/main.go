package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"uuid"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/PuerkitoBio/goquery"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	loginURL    = "https://students.fun.ac.jp/Login"
	booksAPIURL = "https://students.fun.ac.jp/Calendar/Api/Books.ashx"
	targetYear  = "2026"
	targetTerm  = "11"

	userIDSecretName       = "projects/dotto/secrets/FUN_DODAI_CAMPUS_USER_ID/versions/latest"
	userPasswordSecretName = "projects/dotto/secrets/FUN_DODAI_CAMPUS_USER_PASSWORD/versions/latest"
)

// resourceIDToRoomName は dodai campus の resourceId → 教室名 のマッピング。
var resourceIDToRoomName = map[string]string{
	"50": "アトリエ",
	"16": "363",
	"17": "364",
	"18": "365",
	"51": "体育館",
	"2":  "大講義室",
	"19": "483",
	"10": "484",
	"3":  "493",
	"8":  "494C&D",
	"9":  "495C&D",
	"1":  "講堂",
	"11": "583",
	"12": "584",
	"13": "585",
	"4":  "593",
	"5":  "594",
	"6":  "595",
	"14": "R781",
	"15": "R782",
	"7":  "R791",
}

type bookingItem struct {
	ResourceID string `json:"resourceId"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Title      string `json:"title"`
}

func accessSecret(ctx context.Context, name string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanagerクライアント作成失敗: %w", err)
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("シークレット取得失敗 (%s): %w", name, err)
	}
	return string(result.Payload.Data), nil
}

func loginSession(ctx context.Context) (*http.Client, error) {
	username, err := accessSecret(ctx, userIDSecretName)
	if err != nil {
		return nil, err
	}
	password, err := accessSecret(ctx, userPasswordSecretName)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar作成失敗: %w", err)
	}
	client := &http.Client{Jar: jar, Timeout: 30 * time.Second}

	resp, err := client.Get(loginURL)
	if err != nil {
		return nil, fmt.Errorf("ログインページ取得失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ログインページ取得エラー: status=%d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTML解析失敗: %w", err)
	}

	payload := url.Values{}
	payload.Set("__LASTFOCUS", "")
	payload.Set("__EVENTTARGET", "")
	payload.Set("__EVENTARGUMENT", "")
	payload.Set("__SCROLLPOSITIONX", "0")
	payload.Set("__SCROLLPOSITIONY", "0")
	payload.Set("ctl00$MainContent$TargetYearList", targetYear)
	payload.Set("ctl00$MainContent$TargetTermList", targetTerm)
	payload.Set("ctl00$MainContent$LoginId", username)
	payload.Set("ctl00$MainContent$LoginPassword", password)
	payload.Set("ctl00$MainContent$LoginButton", "ログイン")

	for _, name := range []string{"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"} {
		if val, exists := doc.Find(fmt.Sprintf(`[name="%s"]`, name)).Attr("value"); exists {
			payload.Set(name, val)
		}
	}

	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ログインPOST失敗: %w", err)
	}
	defer resp2.Body.Close()
	io.Copy(io.Discard, resp2.Body)

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ログインエラー: status=%d", resp2.StatusCode)
	}

	time.Sleep(2 * time.Second)
	return client, nil
}

func fetchBookings(client *http.Client, t string) ([]bookingItem, error) {
	req, err := http.NewRequest(http.MethodGet, booksAPIURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("t", t)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API取得失敗 (t=%s): %w", t, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込み失敗 (t=%s): %w", t, err)
	}

	var items []bookingItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("JSON解析失敗 (t=%s): %w", t, err)
	}
	return items, nil
}

// filterWeeklySchedule は今日を基準に -7日 〜 +21日 の範囲の予約を返す。
// Python実装と同様に JST で日付境界を計算し、APIのUTC時刻と比較する。
func filterWeeklySchedule(items []bookingItem) []bookingItem {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := time.Now().In(jst)
	todayJST := time.Date(nowJST.Year(), nowJST.Month(), nowJST.Day(), 0, 0, 0, 0, jst)
	from := todayJST.AddDate(0, 0, -7)
	to := todayJST.AddDate(0, 0, 21)

	var result []bookingItem
	for _, item := range items {
		end, err := time.Parse("2006-01-02T15:04:05Z", item.End)
		if err != nil {
			log.Printf("end時刻パース失敗 (resourceId=%s): %v", item.ResourceID, err)
			continue
		}
		endJST := end.In(jst)
		if !endJST.Before(from) && !endJST.After(to) {
			result = append(result, item)
		}
	}
	return result
}

func saveReservations(ctx context.Context, gdb *gorm.DB, items []bookingItem) error {
	var rooms []model.Room
	if err := gdb.WithContext(ctx).Find(&rooms).Error; err != nil {
		return fmt.Errorf("ルーム取得失敗: %w", err)
	}
	roomNameToID := make(map[string]uuid.UUID, len(rooms))
	for _, r := range rooms {
		roomNameToID[r.Name] = r.ID
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := time.Now().In(jst)
	todayJST := time.Date(nowJST.Year(), nowJST.Month(), nowJST.Day(), 0, 0, 0, 0, jst)
	from := todayJST.AddDate(0, 0, -7).UTC()
	to := todayJST.AddDate(0, 0, 21).UTC()

	if err := gdb.WithContext(ctx).
		Where("end_time >= ? AND end_time <= ?", from, to).
		Delete(&model.RoomReservation{}).Error; err != nil {
		return fmt.Errorf("既存予約削除失敗: %w", err)
	}

	var reservations []model.RoomReservation
	for _, item := range items {
		roomName, ok := resourceIDToRoomName[item.ResourceID]
		if !ok {
			continue
		}
		roomID, ok := roomNameToID[roomName]
		if !ok {
			log.Printf("Room not found in DB: %s (resourceId=%s)", roomName, item.ResourceID)
			continue
		}

		startTime, err := time.Parse("2006-01-02T15:04:05Z", item.Start)
		if err != nil {
			log.Printf("start時刻パース失敗 (resourceId=%s): %v", item.ResourceID, err)
			continue
		}
		endTime, err := time.Parse("2006-01-02T15:04:05Z", item.End)
		if err != nil {
			log.Printf("end時刻パース失敗 (resourceId=%s): %v", item.ResourceID, err)
			continue
		}

		reservations = append(reservations, model.RoomReservation{
			Title:     item.Title,
			RoomID:    roomID,
			StartTime: startTime.UTC(),
			EndTime:   endTime.UTC(),
		})
	}

	if len(reservations) == 0 {
		log.Println("保存対象の予約なし")
		return nil
	}

	return gdb.WithContext(ctx).Create(&reservations).Error
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env読み込み失敗(環境変数が既に設定されていれば無視してOK):", err)
	}

	ctx := context.Background()

	client, err := loginSession(ctx)
	if err != nil {
		log.Fatalf("ログイン失敗: %v", err)
	}
	log.Println("ログイン成功")

	var allItems []bookingItem
	for _, t := range []string{"tt", "fr", "sl"} {
		items, err := fetchBookings(client, t)
		if err != nil {
			log.Fatalf("予約取得失敗: %v", err)
		}
		allItems = append(allItems, items...)
		log.Printf("取得完了: t=%s, %d件", t, len(items))
	}

	weekly := filterWeeklySchedule(allItems)
	log.Printf("週間フィルタ後: %d件", len(weekly))

	conn, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("DB切断失敗: %v", err)
		}
	}()

	if err := saveReservations(ctx, conn, weekly); err != nil {
		log.Fatalf("予約保存失敗: %v", err)
	}
	log.Printf("予約保存完了: %d件", len(weekly))
}
