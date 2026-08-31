package handler

import (
	"time"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

func toAPIBusTrip(date time.Time, detail domain.TripDetail) api.BusTrip {
	stops := make([]api.BusStop, 0, len(detail.StopTimes))
	for _, stopTime := range detail.StopTimes {
		stops = append(stops, api.BusStop{
			Id:   stopTime.Stop.StopID,
			Name: stopTime.Stop.StopName,
		})
	}

	var departureTime, arrivalTime time.Time
	if len(detail.StopTimes) > 0 {
		departureTime, _ = gtfsTimeToDate(date, detail.StopTimes[0].StopTime.DepartureTime)
		arrivalTime, _ = gtfsTimeToDate(date, detail.StopTimes[len(detail.StopTimes)-1].StopTime.ArrivalTime)
	}

	return api.BusTrip{
		Id:            detail.Trip.TripID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Route: api.BusRoute{
			Id:   detail.Route.RouteID,
			Name: detail.Route.RouteShortName,
		},
		Stops: stops,
		Delay: "",
		Alert: api.None,
	}
}

func toAPIBusTimetableStop(tripID string, stopTime domain.StopTimeWithStop) api.BusTimetableStop {
	departureTime, _ := parseStopTime(stopTime.StopTime.DepartureTime)
	return api.BusTimetableStop{
		TripId: tripID,
		Stop: api.BusStop{
			Id:   stopTime.Stop.StopID,
			Name: stopTime.Stop.StopName,
		},
		DepartureTime: departureTime,
	}
}

func gtfsTimeToDate(date time.Time, raw string) (time.Time, error) {
	parsed, err := parseStopTime(raw)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, time.UTC), nil
}

func parseStopTime(raw string) (time.Time, error) {
	return time.Parse("15:04:05", raw)
}
