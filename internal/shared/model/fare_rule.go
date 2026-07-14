package model

type FareRule struct {
	Common

	FareID        string  `gorm:"not null;index"`
	RouteID       string  `gorm:"not null;index"`
	Route         *Route  `gorm:"foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OriginID      string  `gorm:"not null;index"`
	Origin        *Stop   `gorm:"foreignKey:OriginID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DestinationID string  `gorm:"not null;index"`
	Destination   *Stop   `gorm:"foreignKey:DestinationID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64 `gorm:"not null"`
}
