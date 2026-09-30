package model

import (
	"github.com/shopspring/decimal"
	"uuid"
)
type RiderCategory string

const (
	RiderCategoryAdult RiderCategory = "adult"
)

// FarePrice は FareRule (どの区間か) に対する、運賃区分ごとの金額を持つ。
// 現状は 1 FareRule につき adult の 1 行のみが作られる想定。
type FarePrice struct {
	Common

	FareRuleID    uuid.UUID       `gorm:"type:uuid;not null;index:idx_fare_price_rule_category,unique"`
	FareRule      *FareRule       `gorm:"belongsTo;foreignKey:FareRuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RiderCategory RiderCategory   `gorm:"type:text;not null;default:'adult';index:idx_fare_price_rule_category,unique"`
	Price         decimal.Decimal `gorm:"type:numeric;not null"`
}
