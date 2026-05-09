package model

import "time"

type Announcement struct {
	Common

	Title          string     `gorm:"not null"`
	URL            string     `gorm:"not null"`
	AvailableFrom  time.Time  `gorm:"not null;index"`
	AvailableUntil *time.Time `gorm:"index"`
}
