package model

// FareRule は「どの経路のどの区間か」だけを表す。実際の金額は FarePrice 側が持つ
// (将来、大人/子ども等の運賃区分が増えても FareRule 自体は変更不要にするため)。
type FareRule struct {
	FareID        string  `gorm:"primaryKey"`
	RouteID       *string `gorm:"index"`
	OriginID      *string `gorm:"index"`
	DestinationID *string `gorm:"index"`
	Route         *Route  `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE"`
}
