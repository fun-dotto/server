package model

type Route struct {
	RouteID        string `gorm:"primaryKey"`
	RouteShortName string `gorm:"not null"`
}
