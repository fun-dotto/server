package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/bus"
)

func (h *Handler) BusTimetableStopsV1List(ctx context.Context, request api.BusTimetableStopsV1ListRequestObject) (api.BusTimetableStopsV1ListResponseObject, error) {
	stopTimes, err := h.busService.ListTimetableStops(ctx, request.TripId)
	if err != nil {
		return nil, err
	}

	apiStops := make([]api.BusTimetableStop, 0, len(stopTimes))
	for _, stopTime := range stopTimes {
		stop, err := h.busService.GetStopByID(ctx, stopTime.StopID)
		if err != nil {
			return nil, err
		}
		apiStops = append(apiStops, toAPIBusTimetableStop(request.TripId, stopTime, stop))
	}

	return api.BusTimetableStopsV1List200JSONResponse{BusTimetableStops: apiStops}, nil
}
