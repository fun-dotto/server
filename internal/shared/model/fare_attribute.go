package model

import (
	"github.com/shopspring/decimal"
	"uuid"
)

type RiderCategory string

const (
	RiderCategoryAdult RiderCategory = "adult"
)

type FareAttribute struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FareID        string          `gorm:"not null;uniqueIndex:idx_fare_rider"`
	RiderCategory RiderCategory   `gorm:"type:text;not null;default:'adult';uniqueIndex:idx_fare_rider"`
	Price         decimal.Decimal `gorm:"type:numeric;not null"`
}
