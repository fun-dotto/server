package model

type FareRule struct {
	Common

	RouteID       string  `gorm:"not null;index"`
	Route         *Route  `gorm:"belongsTo;foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OriginID      string  `gorm:"not null;index"`
	Origin        *Stop   `gorm:"belongsTo;foreignKey:OriginID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DestinationID string  `gorm:"not null;index"`
	Destination   *Stop   `gorm:"belongsTo;foreignKey:DestinationID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64 `gorm:"not null"`
}
