package model

type Stop struct {
	StopID   string  `gorm:"primaryKey"`
	StopName string  `gorm:"not null"`
	StopLat  float64 `gorm:"not null"`
	StopLon  float64 `gorm:"not null"`
}
