package model

type Trip struct {
	TripID      string   `gorm:"primaryKey"`
	RouteID     string   `gorm:"not null;index"`
	ServiceID   string   `gorm:"not null;index"`
	DirectionID *int     `gorm:"check:direction_id IN (0,1)"`
	Route       Route    `gorm:"belongsTo;foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	Calendar    Calendar `gorm:"belongsTo;foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
}
