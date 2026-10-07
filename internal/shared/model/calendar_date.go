package model

import "time"

type CalendarDate struct {
	ServiceID     string    `gorm:"primaryKey"`
	Calendar      *Calendar `gorm:"belongsTo;foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	Date          time.Time `gorm:"type:date;primaryKey"`
	ExceptionType int       `gorm:"not null;check:exception_type IN (1,2)"`
}
