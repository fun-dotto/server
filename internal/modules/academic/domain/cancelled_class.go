package domain

import "time"

type CancelledClass struct {
	ID      string
	Subject Subject
	Date    string
	Period  Period
	Comment string
}

type CancelledClassListFilter struct {
	SubjectIDs []string
	From       *time.Time
	Until      *time.Time
}
