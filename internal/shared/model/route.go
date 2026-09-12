package model

type Route struct {
	Common

	RouteID        string `gorm:"not null;uniqueIndex"`
	RouteShortName string `gorm:"not null"`
}
