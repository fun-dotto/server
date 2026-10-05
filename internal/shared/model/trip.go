package model

type Trip struct {
	TripID      string   `gorm:"primaryKey"`
	RouteID     string   `gorm:"not null;index"`
	Route       *Route   `gorm:"belongsTo;foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	ServiceID   string   `gorm:"not null;index"`
	Calendar    *Calendar`gorm:"belongsTo;foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	DirectionID *int     `gorm:"check:direction_id IN (0,1)"`
}
