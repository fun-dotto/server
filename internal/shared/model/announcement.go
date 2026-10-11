package model

import "time"

type Announcement struct {
	Common

	Title string `gorm:"not null"`
	URL   string `gorm:"not null"`
	// TODO: AIP-140 に従い前置詞を含まない名前（例: available_start_time）に改名する
	AvailableFrom time.Time `gorm:"not null;index"`
	// TODO: AIP-140 に従い前置詞を含まない名前（例: available_end_time）に改名する
	AvailableUntil *time.Time `gorm:"index"`
}
