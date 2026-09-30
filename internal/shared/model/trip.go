package model

type Trip struct {
	TripID      string `gorm:"primaryKey"`
	RouteID     string `gorm:"not null;index"`
	Route       *Route `gorm:"belongsTo;foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServiceID   string `gorm:"not null;index"`
	DirectionID *int
}
