package model

type FareRule struct {
	Common

	RouteID       *string `gorm:"index"`
	Route         *Route  `gorm:"belongsTo;foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OriginID      *string `gorm:"index"`
	Origin        *Stop   `gorm:"belongsTo;foreignKey:OriginID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DestinationID *string `gorm:"index"`
	Destination   *Stop   `gorm:"belongsTo;foreignKey:DestinationID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64 `gorm:"not null"`
}
