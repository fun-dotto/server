package model

import (
	"github.com/shopspring/decimal"
	"uuid"
)
type RiderCategory string

const (
	RiderCategoryAdult RiderCategory = "adult"
)

type FarePrice struct {
	Common

	FareRuleID    uuid.UUID       `gorm:"type:uuid;not null;index:idx_fare_price_rule_category,unique"`
	FareRule      *FareRule       `gorm:"belongsTo;foreignKey:FareRuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RiderCategory RiderCategory   `gorm:"type:text;not null;default:'adult';index:idx_fare_price_rule_category,unique"`
	Price         decimal.Decimal `gorm:"type:numeric;not null"`
}
