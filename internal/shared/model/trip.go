package model

type Trip struct {
	TripID      string   `gorm:"primaryKey"`
	RouteID     string   `gorm:"not null;index"`
	ServiceID   string   `gorm:"not null;index"`
	DirectionID *int     `gorm:"check:direction_id IN (0,1)"`
	Route       Route    `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE"`
	Calendar    Calendar `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
}
