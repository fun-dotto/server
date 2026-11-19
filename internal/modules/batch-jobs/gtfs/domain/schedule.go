package domain

import "github.com/fun-dotto/server/internal/shared/model"

type Schedule struct {
	Calendars      []model.Calendar
	CalendarDates  []model.CalendarDate
	Zones          []model.Zone
	FareAttributes []model.FareAttribute
	FareRules      []model.FareRule
	Routes         []model.Route
	Stops          []model.Stop
	Trips          []model.Trip
	StopTimes      []model.StopTime
}
