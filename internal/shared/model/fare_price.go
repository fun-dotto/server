package model

import (
	"github.com/shopspring/decimal"
)

type RiderCategory string

const (
	RiderCategoryAdult RiderCategory = "adult"
)

type FarePrice struct {
	FareID        string          `gorm:"primaryKey"`
	RiderCategory RiderCategory   `gorm:"type:text;primaryKey;default:'adult'"`
	Price         decimal.Decimal `gorm:"type:numeric;not null"`
	FareRule      FareRule `gorm:"foreignKey:FareID;constraint:OnDelete:CASCADE"`
}
