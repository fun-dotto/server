package model

type CalendarDate struct {
	Common

	ServiceID     string `gorm:"not null;index"`
	Date          string `gorm:"not null;index"`
	ExceptionType int    `gorm:"not null"`
}
