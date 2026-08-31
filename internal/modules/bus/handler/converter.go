package handler

import (
	"time"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

func toAPIBusTrip(date time.Time, trip domain.Trip, route domain.Route, stopTimes []domain.StopTime, stopByID map[string]domain.Stop) api.BusTrip {
	stops := make([]api.BusStop, 0, len(stopTimes))
	for _, stopTime := range stopTimes {
		stop, ok := stopByID[stopTime.StopID]
		if !ok {
			continue
		}
		stops = append(stops, api.BusStop{Id: stop.StopID, Name: stop.StopName})
	}

	departureTime, _ := gtfsTimeToDate(date, stopTimes[0].DepartureTime)
	arrivalTime, _ := gtfsTimeToDate(date, stopTimes[len(stopTimes)-1].ArrivalTime)

	return api.BusTrip{
		Id:            trip.TripID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Route: api.BusRoute{
			Id:   route.RouteID,
			Name: route.RouteShortName,
		},
		Stops: stops,
		Delay: "",
		Alert: api.None,
	}
}

func toAPIBusTimetableStop(tripID string, stopTime domain.StopTime, stop domain.Stop) api.BusTimetableStop {
	departureTime, _ := parseStopTime(stopTime.DepartureTime)
	return api.BusTimetableStop{
		TripId:        tripID,
		Stop:          api.BusStop{Id: stop.StopID, Name: stop.StopName},
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
