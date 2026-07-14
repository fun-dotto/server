package model

type Trip struct {
	Common

	TripID      string    `gorm:"not null;uniqueIndex"`
	RouteID     string    `gorm:"not null;index"`
	Route       *Route    `gorm:"foreignKey:RouteID;references:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServiceID   string    `gorm:"not null;index"`
	Calendar    *Calendar `gorm:"foreignKey:ServiceID;references:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DirectionID int       `gorm:"not null"`
}
