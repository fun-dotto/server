package model

type Stop struct {
	Common

	StopID   string `gorm:"not null;uniqueIndex"`
	StopName string `gorm:"not null"`
}
