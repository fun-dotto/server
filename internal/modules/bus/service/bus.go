package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

type BusRepository interface {
	ListTripsByDate(ctx context.Context, date time.Time) ([]domain.Trip, error)
	GetRoutesByIDs(ctx context.Context, routeIDs []string) (map[string]domain.Route, error)
	ListStopTimesByTripIDs(ctx context.Context, tripIDs []string) ([]domain.StopTime, error)
	GetStopsByIDs(ctx context.Context, stopIDs []string) (map[string]domain.Stop, error)
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
	if len(trips) == 0 {
		return []domain.TripDetail{}, nil
	}

	routeIDs := uniqueIDs(len(trips), func(i int) string { return trips[i].RouteID })
	routes, err := s.busRepository.GetRoutesByIDs(ctx, routeIDs)
	if err != nil {
		return nil, err
	}

	tripIDs := make([]string, len(trips))
	for i, trip := range trips {
		tripIDs[i] = trip.TripID
	}

	stopTimes, err := s.busRepository.ListStopTimesByTripIDs(ctx, tripIDs)
	if err != nil {
		return nil, err
	}

	stopIDs := uniqueIDs(len(stopTimes), func(i int) string { return stopTimes[i].StopID })
	stops, err := s.busRepository.GetStopsByIDs(ctx, stopIDs)
	if err != nil {
		return nil, err
	}

	stopTimesByTripID := make(map[string][]domain.StopTimeWithStop, len(trips))
	for _, stopTime := range stopTimes {
		stop, ok := stops[stopTime.StopID]
		if !ok {
			return nil, fmt.Errorf("stop not found: %s", stopTime.StopID)
		}
		stopTimesByTripID[stopTime.TripID] = append(stopTimesByTripID[stopTime.TripID], domain.StopTimeWithStop{
			StopTime: stopTime,
			Stop:     stop,
		})
	}

	details := make([]domain.TripDetail, 0, len(trips))
	for _, trip := range trips {
		route, ok := routes[trip.RouteID]
		if !ok {
			return nil, fmt.Errorf("route not found: %s", trip.RouteID)
		}
		details = append(details, domain.TripDetail{
			Trip:      trip,
			Route:     route,
			StopTimes: stopTimesByTripID[trip.TripID],
		})
	}
	return details, nil
}

func (s *BusService) ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTimeWithStop, error) {
	stopTimes, err := s.busRepository.ListStopTimesByTripIDs(ctx, []string{tripID})
	if err != nil {
		return nil, err
	}

	stopIDs := uniqueIDs(len(stopTimes), func(i int) string { return stopTimes[i].StopID })
	stops, err := s.busRepository.GetStopsByIDs(ctx, stopIDs)
	if err != nil {
		return nil, err
	}

	out := make([]domain.StopTimeWithStop, 0, len(stopTimes))
	for _, stopTime := range stopTimes {
		stop, ok := stops[stopTime.StopID]
		if !ok {
			return nil, fmt.Errorf("stop not found: %s", stopTime.StopID)
		}
		out = append(out, domain.StopTimeWithStop{
			StopTime: stopTime,
			Stop:     stop,
		})
	}
	return out, nil
}

func uniqueIDs(n int, at func(i int) string) []string {
	seen := make(map[string]struct{}, n)
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		id := at(i)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
