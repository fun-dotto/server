package model

type FareRule struct {
	FareID        string  `gorm:"primaryKey"`
	RouteID       *string `gorm:"index"`
	OriginID      *string `gorm:"index"`
	DestinationID *string `gorm:"index"`
	Route         Route   `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE"`
}
