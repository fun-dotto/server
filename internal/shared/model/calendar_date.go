package model

type CalendarDate struct {
	Common

	ServiceID     string    `gorm:"not null;index"`
	Calendar      *Calendar `gorm:"belongsTo;foreignKey:ServiceID;references:ServiceID"`
	Date          string    `gorm:"not null;index"`
	ExceptionType int       `gorm:"not null"`
}
