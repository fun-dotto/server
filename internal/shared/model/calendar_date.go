package model

type CalendarDate struct {
	Common

	ServiceID     string    `gorm:"not null;index"`
	Calendar      *Calendar `gorm:"belongsTo;foreignKey:ServiceID;references:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Date          string    `gorm:"not null;index"`
	ExceptionType int       `gorm:"not null"`
}
