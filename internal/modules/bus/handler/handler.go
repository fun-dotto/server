package handler

import (
	"context"
	"time"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type busService interface {
	ListTrips(ctx context.Context, date time.Time) ([]domain.Trip, error)
	GetRouteByID(ctx context.Context, routeID string) (domain.Route, error)
	ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTime, error)
	GetStopByID(ctx context.Context, stopID string) (domain.Stop, error)
}

type Handler struct {
	busService busService
}

func NewHandler(busService busService) *Handler {
	return &Handler{busService: busService}
}
