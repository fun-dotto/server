package service

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

type BusRepository interface {
	ListTripsByDate(ctx context.Context, date time.Time) ([]domain.Trip, error)
	GetRouteByID(ctx context.Context, routeID string) (domain.Route, error)
	ListStopTimesByTripID(ctx context.Context, tripID string) ([]domain.StopTime, error)
	GetStopByID(ctx context.Context, stopID string) (domain.Stop, error)
}

type BusService struct {
	busRepository BusRepository
}

func NewBusService(busRepository BusRepository) *BusService {
	return &BusService{busRepository: busRepository}
}

func (s *BusService) ListTrips(ctx context.Context, date time.Time) ([]domain.Trip, error) {
	return s.busRepository.ListTripsByDate(ctx, date)
}

func (s *BusService) GetRouteByID(ctx context.Context, routeID string) (domain.Route, error) {
	return s.busRepository.GetRouteByID(ctx, routeID)
}

func (s *BusService) ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTime, error) {
	return s.busRepository.ListStopTimesByTripID(ctx, tripID)
}

func (s *BusService) GetStopByID(ctx context.Context, stopID string) (domain.Stop, error) {
	return s.busRepository.GetStopByID(ctx, stopID)
}
