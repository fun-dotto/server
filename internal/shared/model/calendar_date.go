package model

import "time"

type CalendarDate struct {
	ServiceID     string    `gorm:"primaryKey"`
	Date          time.Time `gorm:"type:date;primaryKey"`
	ExceptionType bool      `gorm:"not null"`
	Calendar      Calendar  `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
}
