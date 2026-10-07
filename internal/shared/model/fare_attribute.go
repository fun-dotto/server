package model

import "github.com/shopspring/decimal"

type RiderCategory string

const (
	RiderCategoryAdult RiderCategory = "adult"
)

type FareAttribute struct {
	FareID        string          `gorm:"primaryKey"`
	RiderCategory RiderCategory   `gorm:"type:text;not null;default:'adult'"`
	Price         decimal.Decimal `gorm:"type:numeric;not null"`
}
