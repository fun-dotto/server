package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/bus"
)

func (h *Handler) BusTripsV1List(ctx context.Context, request api.BusTripsV1ListRequestObject) (api.BusTripsV1ListResponseObject, error) {
	date := request.Params.Date.Time
	details, err := h.busService.ListTripDetails(ctx, date)
	if err != nil {
		return nil, err
	}

	apiTrips := make([]api.BusTrip, 0, len(details))
	for _, detail := range details {
		apiTrips = append(apiTrips, toAPIBusTrip(date, detail))
	}

	return api.BusTripsV1List200JSONResponse{BusTrips: apiTrips}, nil
}
