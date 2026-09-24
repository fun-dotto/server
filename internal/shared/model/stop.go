package model

type Stop struct {
	Common

	StopID   string  `gorm:"not null;uniqueIndex"`
	StopName string  `gorm:"not null"`
	StopLat  float64 `gorm:"not null"`
	StopLon  float64 `gorm:"not null"`
}
