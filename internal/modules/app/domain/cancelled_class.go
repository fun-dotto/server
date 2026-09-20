package domain

import "time"

// CancelledClass 休講
type CancelledClass struct {
	ID      string
	Comment string
	Date    time.Time
	Period  Period
	Subject Subject
}

// CancelledClassQuery 休講検索クエリ
type CancelledClassQuery struct {
	SubjectIDs *[]string
	From       *time.Time
	Until      *time.Time
}
