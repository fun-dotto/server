package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

const (
	loginURL   = "https://students.fun.ac.jp/Login"
	targetYear = "2026" // 実際の年度に置き換えてください
	targetTerm = "1"

	targetAPIURL = "https://students.fun.ac.jp/Calendar/Api/Books.ashx?t=tt"
)

// loginSession はポータルにログインし、Cookie(セッション)を保持したhttp.Clientを返す
func loginSession() (*http.Client, error) {
	username := os.Getenv("USER_ID")
	password := os.Getenv("USER_PASSWORD")
	if username == "" || password == "" {
		return nil, fmt.Errorf("環境変数 USER_ID / USER_PASSWORD が設定されていません")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar作成失敗: %w", err)
	}
	client := &http.Client{Jar: jar}

	// --- 1. ログインページをGETして隠しフィールドを取得 ---
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

	// --- 2. フォームのpayloadを組み立てる ---
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

	hiddenFields := []string{"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"}
	for _, name := range hiddenFields {
		selector := fmt.Sprintf(`[name="%s"]`, name)
		if val, exists := doc.Find(selector).Attr("value"); exists {
			payload.Set(name, val)
		}
	}

	// --- 3. フォームをPOSTしてログイン ---
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

	fmt.Println("POST後の実際のURL:", resp2.Request.URL.String())

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ログインエラー: status=%d", resp2.StatusCode)
	}

	time.Sleep(2 * time.Second)

	return client, nil
}

// fetchAndPrint は認証済みclientでAPIを叩いて中身を表示する
func fetchAndPrint(client *http.Client, apiURL string) error {
	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("API取得失敗: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("ステータスコード:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("読み込み失敗: %w", err)
	}

	fmt.Println("---レスポンス内容---")
	fmt.Println(string(body))

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env読み込み失敗(環境変数が既に設定されていれば無視してOK):", err)
	}

	// 1. ログインしてセッション付きclientを取得
	client, err := loginSession()
	if err != nil {
		log.Fatalf("ログイン失敗: %v", err)
	}
	fmt.Println("ログイン成功。")

	// 2. そのclientでAPIを叩いて中身を表示
	if err := fetchAndPrint(client, targetAPIURL); err != nil {
		log.Fatalf("API取得失敗: %v", err)
	}
}
