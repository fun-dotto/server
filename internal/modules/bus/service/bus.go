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

func (s *BusService) ListTripDetails(ctx context.Context, date time.Time) ([]domain.TripDetail, error) {
	trips, err := s.busRepository.ListTripsByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	details := make([]domain.TripDetail, 0, len(trips))
	for _, trip := range trips {
		route, err := s.busRepository.GetRouteByID(ctx, trip.RouteID)
		if err != nil {
			return nil, err
		}
		stopTimes, err := s.listStopTimesWithStops(ctx, trip.TripID)
		if err != nil {
			return nil, err
		}
		details = append(details, domain.TripDetail{
			Trip:      trip,
			Route:     route,
			StopTimes: stopTimes,
		})
	}
	return details, nil
}

func (s *BusService) ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTimeWithStop, error) {
	return s.listStopTimesWithStops(ctx, tripID)
}

func (s *BusService) listStopTimesWithStops(ctx context.Context, tripID string) ([]domain.StopTimeWithStop, error) {
	stopTimes, err := s.busRepository.ListStopTimesByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	out := make([]domain.StopTimeWithStop, 0, len(stopTimes))
	for _, stopTime := range stopTimes {
		stop, err := s.busRepository.GetStopByID(ctx, stopTime.StopID)
		if err != nil {
			return nil, err
		}
		out = append(out, domain.StopTimeWithStop{
			StopTime: stopTime,
			Stop:     stop,
		})
	}
	return out, nil
}
