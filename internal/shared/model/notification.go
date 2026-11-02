package model

import "time"

type Notification struct {
	// Deprecated: Use Common instead
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`

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

	// TODO: AIP-140 に従い前置詞を含まない名前（例: notify_start_time）に改名する
	NotifyAfter time.Time `gorm:"not null;index"`
	// TODO: AIP-140 に従い前置詞を含まない名前（例: notify_end_time）に改名する
	NotifyBefore time.Time `gorm:"not null;index"`
}
