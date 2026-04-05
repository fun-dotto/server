package domain

import "time"

type MakeupClass struct {
	ID      string
	Subject Subject
	Date    string
	Period  Period
	Comment string
}

type MakeupClassListFilter struct {
	SubjectIDs []string
	From       *time.Time
	Until      *time.Time
}
