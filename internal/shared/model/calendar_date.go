package model

import "time"

type CalendarDate struct {
	ServiceID     string    `gorm:"primaryKey"`
	Date          time.Time `gorm:"type:date;primaryKey"`
	ExceptionType int       `gorm:"not null;check:exception_type IN (1,2)"`
}
