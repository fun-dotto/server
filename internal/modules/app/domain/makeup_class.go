package domain

import "time"

// MakeupClass 補講
type MakeupClass struct {
	ID      string
	Comment string
	Date    time.Time
	Period  Period
	Subject Subject
}

// MakeupClassQuery 補講検索クエリ
type MakeupClassQuery struct {
	SubjectIDs *[]string
	From       *time.Time
	Until      *time.Time
}
