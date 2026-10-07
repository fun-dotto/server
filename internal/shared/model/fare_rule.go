package model

import "uuid"

type FareRule struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FareID        string         `gorm:"not null;index"`
	FareAttribute *FareAttribute `gorm:"belongsTo;foreignKey:FareID;constraint:OnDelete:CASCADE"`
	RouteID       *string        `gorm:"index"`
	Route         *Route         `gorm:"belongsTo;foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	OriginID      *string        `gorm:"index"`
	Origin        *Zone          `gorm:"foreignKey:OriginID;constraint:OnDelete:CASCADE"`
	DestinationID *string        `gorm:"index"`
	Destination   *Zone          `gorm:"foreignKey:DestinationID;constraint:OnDelete:CASCADE"`
}
