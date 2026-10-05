package model

type FareRule struct {
	FareID        string  `gorm:"primaryKey"`
	RouteID       *string `gorm:"index"`
	Route         *Route  `gorm:"belongsTo;foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	OriginID      *string `gorm:"index"`
	Origin        *Zone   `gorm:"foreignKey:OriginID;constraint:OnDelete:CASCADE"`
	DestinationID *string `gorm:"index"`
	Destination   *Zone   `gorm:"foreignKey:DestinationID;constraint:OnDelete:CASCADE"`
}
