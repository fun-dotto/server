package handler

import (
	"context"
	"time"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type busService interface {
	ListTripDetails(ctx context.Context, date time.Time) ([]domain.TripDetail, error)
	ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTimeWithStop, error)
}

type Handler struct {
	busService busService
}

func NewHandler(busService busService) *Handler {
	return &Handler{busService: busService}
}
