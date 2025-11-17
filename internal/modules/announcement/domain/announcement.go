package domain

import "time"

type Announcement struct {
	ID       string
	Title    string
	Date     time.Time
	URL      string
	IsActive bool
}
