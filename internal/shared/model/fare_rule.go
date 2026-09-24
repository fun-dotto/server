package model

// FareRule は「どの経路のどの区間か」だけを表す。実際の金額は FarePrice 側が持つ
// (将来、大人/子ども等の運賃区分が増えても FareRule 自体は変更不要にするため)。
type FareRule struct {
	Common

	RouteID       *string `gorm:"index"`
	Route         *Route  `gorm:"belongsTo;foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OriginID      *string `gorm:"index"`
	Origin        *Stop   `gorm:"belongsTo;foreignKey:OriginID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DestinationID *string `gorm:"index"`
	Destination   *Stop   `gorm:"belongsTo;foreignKey:DestinationID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// FarePrice 側の belongsTo だけに constraint を書くと GORM が関係定義を上書きし
	// FK の CASCADE が効かなくなるため、hasMany 側にも同じ constraint を明示している。
	Prices []FarePrice `gorm:"hasMany;foreignKey:FareRuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
