package service

import (
	"github.com/fun-dotto/server/internal/modules/bus/domain"
	"github.com/fun-dotto/server/internal/modules/bus/gtfs"
)

// ConvertToDomain は GTFS の各ファイルを読み、domain に変換する。
func (s *ScheduleService) ConvertToDomain(archive *gtfs.Archive) (
	stops []domain.Stop,
	routes []domain.Route,
	calendars []domain.Calendar,
	calendarDates []domain.CalendarDate,
	trips []domain.Trip,
	stopTimes []domain.StopTime,
	fareRules []domain.FareRule,
	err error,
) {
	rawStops, err := archive.Stops()
	if err != nil {
		return
	}
	if stops, err = gtfs.ToDomainStops(rawStops); err != nil {
		return
	}

	rawRoutes, err := archive.Routes()
	if err != nil {
		return
	}
	if routes, err = gtfs.ToDomainRoutes(rawRoutes); err != nil {
		return
	}

	rawCalendars, err := archive.Calendars()
	if err != nil {
		return
	}
	if calendars, err = gtfs.ToDomainCalendars(rawCalendars); err != nil {
		return
	}

	rawCalendarDates, err := archive.CalendarDates()
	if err != nil {
		return
	}
	if calendarDates, err = gtfs.ToDomainCalendarDates(rawCalendarDates); err != nil {
		return
	}

	rawTrips, err := archive.Trips()
	if err != nil {
		return
	}
	if trips, err = gtfs.ToDomainTrips(rawTrips); err != nil {
		return
	}

	rawStopTimes, err := archive.StopTimes()
	if err != nil {
		return
	}
	if stopTimes, err = gtfs.ToDomainStopTimes(rawStopTimes); err != nil {
		return
	}

	rawFareRules, rawFareAttrs, err := archive.FareTables()
	if err != nil {
		return
	}
	fareRules, err = gtfs.ToDomainFareRules(rawFareRules, rawFareAttrs)
	return
}
