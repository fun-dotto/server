package model

import "time"

type CalendarDate struct {
	Common

	ServiceID     string    `gorm:"not null;index"`
	Date          time.Time `gorm:"type:date;not null;index"`
	ExceptionType int       `gorm:"not null"`
}
