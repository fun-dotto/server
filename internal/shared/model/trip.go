package model

type Trip struct {
	Common

	TripID      string `gorm:"not null;uniqueIndex"`
	RouteID     string `gorm:"not null;index"`
	Route       *Route `gorm:"belongsTo;foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServiceID   string `gorm:"not null;index"`
	DirectionID *int
}
