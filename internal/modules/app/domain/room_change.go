package domain

import "time"

// RoomChange 教室変更
type RoomChange struct {
	ID           string
	Date         time.Time
	Period       Period
	Subject      Subject
	OriginalRoom Room
	NewRoom      Room
}

// RoomChangeQuery 教室変更検索クエリ
type RoomChangeQuery struct {
	SubjectIDs *[]string
	From       *time.Time
	Until      *time.Time
}
