package model

import "time"

type CalendarDate struct {
	Common

	ServiceID     string    `gorm:"not null;uniqueIndex:idx_calendar_dates_service_date"`
	Date          time.Time `gorm:"type:date;not null;index;uniqueIndex:idx_calendar_dates_service_date"`
	ExceptionType int       `gorm:"not null;check:exception_type IN (1,2)"`
}
