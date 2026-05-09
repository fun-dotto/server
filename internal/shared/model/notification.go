package model

import "time"

type Notification struct {
	Common

	Title                string `gorm:"not null"`
	Body                 string `gorm:"not null"`
	ImageURL             *string
	AnalyticsLabel       *string
	APNsBadge            *int
	APNsSound            *string
	APNsContentAvailable *bool
	AndroidChannelID     *string
	AndroidPriority      *string
	AndroidTTLSeconds    *int
	WebpushLink          *string

	URL *string

	NotifyAfter  time.Time `gorm:"not null;index"`
	NotifyBefore time.Time `gorm:"not null;index"`
}
