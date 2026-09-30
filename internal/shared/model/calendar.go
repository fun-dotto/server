package model

import "time"

type Calendar struct {
	ServiceID string    `gorm:"primaryKey"`
	Monday    bool      `gorm:"not null"`
	Tuesday   bool      `gorm:"not null"`
	Wednesday bool      `gorm:"not null"`
	Thursday  bool      `gorm:"not null"`
	Friday    bool      `gorm:"not null"`
	Saturday  bool      `gorm:"not null"`
	Sunday    bool      `gorm:"not null"`
	StartDate time.Time `gorm:"type:date;not null"`
	EndDate   time.Time `gorm:"type:date;not null"`
}
