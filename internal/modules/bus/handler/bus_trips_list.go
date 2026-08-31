package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

func (h *Handler) BusTripsV1List(ctx context.Context, request api.BusTripsV1ListRequestObject) (api.BusTripsV1ListResponseObject, error) {
	date := request.Params.Date.Time
	trips, err := h.busService.ListTrips(ctx, date)
	if err != nil {
		return nil, err
	}

	apiTrips := make([]api.BusTrip, 0, len(trips))
	for _, trip := range trips {
		route, err := h.busService.GetRouteByID(ctx, trip.RouteID)
		if err != nil {
			return nil, err
		}
		stopTimes, err := h.busService.ListTimetableStops(ctx, trip.TripID)
		if err != nil {
			return nil, err
		}
		stopByID := make(map[string]domain.Stop, len(stopTimes))
		for _, stopTime := range stopTimes {
			stop, err := h.busService.GetStopByID(ctx, stopTime.StopID)
			if err != nil {
				return nil, err
			}
			stopByID[stopTime.StopID] = stop
		}
		apiTrips = append(apiTrips, toAPIBusTrip(date, trip, route, stopTimes, stopByID))
	}

	return api.BusTripsV1List200JSONResponse{BusTrips: apiTrips}, nil
}
