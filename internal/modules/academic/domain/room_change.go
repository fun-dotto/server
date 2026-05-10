package domain

import "time"

type RoomChange struct {
	ID           string
	Subject      Subject
	Date         string
	Period       Period
	OriginalRoom Room
	NewRoom      Room
}

type RoomChangeListFilter struct {
	SubjectIDs []string
	From       *time.Time
	Until      *time.Time
}
