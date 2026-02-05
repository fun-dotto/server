package domain

import "time"

type Announcement struct {
	ID             string
	Title          string
	URL            string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
}
