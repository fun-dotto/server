package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

const monthDayLayout = "01/02"

// columnTitle は各テーブルの td[data-col-responsive-title] に対応する列名。
const (
	colDate             = "日付"
	colDayOfWeek        = "曜日"
	colPeriod           = "時限"
	colLessonName       = "授業名"
	colCampus           = "キャンパス"
	colStaff            = "代表教職員"
	colCancelledComment = "休講コメント"
	colMakeupComment    = "補講コメント"
	colRoomName         = "教室名"
	colRoomFrom         = "移動元"
	colRoomTo           = "移動先"
)

// ParseClassChanges は /Pt/CSLecture の HTML から休講・補講・部屋変更を抽出する。
// year は対象年度（4/1 起点）で、"MM/DD" の日付を年度内に収まるよう補正するために使う。
func ParseClassChanges(html string, year int) ([]domain.ScrapedCancelledClass, []domain.ScrapedMakeupClass, []domain.ScrapedRoomChange, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse html: %w", err)
	}

	var cancelled []domain.ScrapedCancelledClass
	var makeup []domain.ScrapedMakeupClass
	var roomChanges []domain.ScrapedRoomChange

	doc.Find("tr").Each(func(_ int, row *goquery.Selection) {
		date, period, lessonName, campus, staff, ok := parseCommonColumns(row, year)
		if !ok {
			return
		}

		if comment, ok := cellText(row, colCancelledComment); ok {
			cancelled = append(cancelled, domain.ScrapedCancelledClass{
				Date:       date,
				Period:     period,
				LessonName: lessonName,
				Campus:     campus,
				Staff:      staff,
				Comment:    comment,
				Kind:       classifyCancelledComment(comment),
			})
			return
		}

		if comment, ok := cellText(row, colMakeupComment); ok {
			roomName, _ := cellText(row, colRoomName)
			makeup = append(makeup, domain.ScrapedMakeupClass{
				Date:       date,
				Period:     period,
				LessonName: lessonName,
				Campus:     campus,
				Staff:      staff,
				Comment:    comment,
				RoomName:   roomName,
			})
			return
		}

		roomFrom, hasFrom := cellText(row, colRoomFrom)
		roomTo, hasTo := cellText(row, colRoomTo)
		if hasFrom && hasTo {
			roomChanges = append(roomChanges, domain.ScrapedRoomChange{
				Date:         date,
				Period:       period,
				LessonName:   lessonName,
				Campus:       campus,
				Staff:        staff,
				RoomFromName: roomFrom,
				RoomToName:   roomTo,
			})
		}
	})

	return cancelled, makeup, roomChanges, nil
}

// parseCommonColumns は 3 種のテーブルに共通する日付・時限・授業名・キャンパス・代表教職員を読む。
// 曜日列の有無を「対象行」かどうかの判定に使う（見出し行やその他の tr を除外する）。
func parseCommonColumns(row *goquery.Selection, year int) (date time.Time, period, lessonName, campus, staff string, ok bool) {
	if _, hasDay := cellText(row, colDayOfWeek); !hasDay {
		return time.Time{}, "", "", "", "", false
	}

	dateText, hasDate := cellText(row, colDate)
	periodText, hasPeriod := cellText(row, colPeriod)
	lessonName, hasLesson := cellText(row, colLessonName)
	campus, hasCampus := cellText(row, colCampus)
	staff, hasStaff := cellText(row, colStaff)
	if !hasDate || !hasPeriod || !hasLesson || !hasCampus || !hasStaff {
		return time.Time{}, "", "", "", "", false
	}

	monthDay, err := time.Parse(monthDayLayout, dateText)
	if err != nil {
		return time.Time{}, "", "", "", "", false
	}
	date = resolveNendoDate(monthDay, year)

	period, ok = formatPeriod(periodText)
	if !ok {
		return time.Time{}, "", "", "", "", false
	}

	return date, period, lessonName, campus, staff, true
}

// cellText は td[data-col-responsive-title=name] のテキストを取得する。列が存在しなければ ok=false。
func cellText(row *goquery.Selection, name string) (string, bool) {
	cell := row.Find(fmt.Sprintf(`td[data-col-responsive-title="%s"]`, name))
	if cell.Length() == 0 {
		return "", false
	}
	return strings.TrimSpace(cell.First().Text()), true
}

// formatPeriod は "1時限" 等の時限表示の先頭の数字から "Period{n}" を作る。
func formatPeriod(s string) (string, bool) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", false
	}
	first := []rune(trimmed)[0]
	if first < '0' || first > '9' {
		return "", false
	}
	n, err := strconv.Atoi(string(first))
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("Period%d", n), true
}

// classifyCancelledComment は休講コメントの接頭辞から補講予定の種別を判定する。
func classifyCancelledComment(comment string) domain.CancelledClassKind {
	switch {
	case strings.HasPrefix(comment, "補講あり"):
		return domain.CancelledClassKindMakeupScheduled
	case strings.HasPrefix(comment, "補講なし"):
		return domain.CancelledClassKindMakeupNotPlanned
	case strings.HasPrefix(comment, "補講未定"):
		return domain.CancelledClassKindMakeupUndecided
	default:
		return domain.CancelledClassKindOther
	}
}
