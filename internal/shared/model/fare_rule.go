package model

type FareRule struct {
	FareID        string         `gorm:"primaryKey"`
	FareAttribute *FareAttribute `gorm:"foreignKey:FareID;constraint:OnDelete:CASCADE"`
	RouteID       *string        `gorm:"primaryKey;default:NULL"`
	Route         *Route         `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	OriginID      *string        `gorm:"primaryKey;default:NULL"`
	Origin        *Zone          `gorm:"foreignKey:OriginID;constraint:OnDelete:CASCADE"`
	DestinationID *string        `gorm:"primaryKey;default:NULL"`
	Destination   *Zone          `gorm:"foreignKey:DestinationID;constraint:OnDelete:CASCADE"`
}
