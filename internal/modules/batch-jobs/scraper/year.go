package scraper

import "time"

// nendoStart は年度（4 月始まり）の開始日を返す。
func nendoStart(year int) time.Time {
	return time.Date(year, time.April, 1, 0, 0, 0, 0, time.UTC)
}

// nendoEnd は年度（4 月始まり）の終了日を返す。
func nendoEnd(year int) time.Time {
	return time.Date(year+1, time.March, 31, 0, 0, 0, 0, time.UTC)
}

// CurrentNendo は現在日時を基準とした年度（4/1 起点）を返す。
func CurrentNendo(now time.Time) int {
	if now.Month() < time.April {
		return now.Year() - 1
	}
	return now.Year()
}

// resolveNendoDate は "MM/DD" 形式の日付をポータルの表示年度 year の範囲
// （nendoStart(year) から nendoEnd(year) まで）に収まるよう年を補正する。
func resolveNendoDate(monthDay time.Time, year int) time.Time {
	candidate := time.Date(year, monthDay.Month(), monthDay.Day(), 0, 0, 0, 0, time.UTC)
	if candidate.Before(nendoStart(year)) {
		candidate = time.Date(year+1, monthDay.Month(), monthDay.Day(), 0, 0, 0, 0, time.UTC)
	}
	if candidate.After(nendoEnd(year)) {
		candidate = time.Date(year-1, monthDay.Month(), monthDay.Day(), 0, 0, 0, 0, time.UTC)
	}
	return candidate
}
