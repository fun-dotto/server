package domain

import "time"

var japanStandardTime = time.FixedZone("JST", 9*3600)

// CurrentAcademicYear は now が属する日本の年度を返す。
// 日本時間で 4月1日〜翌年3月31日を同一年度とし、1〜3月は前暦年を年度として扱う。
func CurrentAcademicYear(now time.Time) int {
	nowJST := now.In(japanStandardTime)
	if nowJST.Month() >= time.April {
		return nowJST.Year()
	}
	return nowJST.Year() - 1
}
