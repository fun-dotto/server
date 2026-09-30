package model

import "github.com/shopspring/decimal"

type Stop struct {
	StopID   string          `gorm:"primaryKey"`
	StopName string          `gorm:"not null"`
	StopLat  decimal.Decimal `gorm:"type:numeric(9,6);not null"`
	StopLon  decimal.Decimal `gorm:"type:numeric(9,6);not null"`
	ZoneID   *string         `gorm:"index"`
}
