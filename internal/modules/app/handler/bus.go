package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/app"
)

func (h *Handler) BusTripsV1List(ctx context.Context, request api.BusTripsV1ListRequestObject) (api.BusTripsV1ListResponseObject, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) BusTimetableStopsV1List(ctx context.Context, request api.BusTimetableStopsV1ListRequestObject) (api.BusTimetableStopsV1ListResponseObject, error) {
	return nil, errors.New("not implemented")
}
