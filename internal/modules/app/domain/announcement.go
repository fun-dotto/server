package domain

import "time"

type Announcement struct {
	ID             string
	Title          string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
	URL            string
}

type AnnouncementRequest struct {
	Title          string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
	URL            string
}
