package repository

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type BusRepository struct {
	db *gorm.DB
}

func NewBusRepository(db *gorm.DB) *BusRepository {
	return &BusRepository{db: db}
}

func (r *BusRepository) ListTripsByDate(ctx context.Context, date time.Time) ([]domain.Trip, error) {
	var trips []model.Trip
	if err := r.db.WithContext(ctx).Preload("Route").Find(&trips).Error; err != nil {
		return nil, err
	}

	out := make([]domain.Trip, 0, len(trips))
	for _, trip := range trips {
		active, err := r.isTripActive(ctx, trip.ServiceID, date)
		if err != nil {
			return nil, err
		}
		if !active {
			continue
		}
		out = append(out, domain.Trip{
			ID:          trip.ID.String(),
			TripID:      trip.TripID,
			RouteID:     trip.RouteID,
			ServiceID:   trip.ServiceID,
			DirectionID: trip.DirectionID,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].TripID < out[j].TripID
	})
	return out, nil
}

func (r *BusRepository) GetRouteByID(ctx context.Context, routeID string) (domain.Route, error) {
	var record model.Route
	if err := r.db.WithContext(ctx).Where("route_id = ?", routeID).First(&record).Error; err != nil {
		return domain.Route{}, err
	}
	return domain.Route{ID: record.ID.String(), RouteID: record.RouteID, RouteShortName: record.RouteShortName}, nil
}

func (r *BusRepository) ListStopTimesByTripID(ctx context.Context, tripID string) ([]domain.StopTime, error) {
	var records []model.StopTime
	if err := r.db.WithContext(ctx).Where("trip_id = ?", tripID).Order("stop_sequence ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	out := make([]domain.StopTime, len(records))
	for i, record := range records {
		out[i] = domain.StopTime{
			ID:            record.ID.String(),
			TripID:        record.TripID,
			ArrivalTime:   record.ArrivalTime,
			DepartureTime: record.DepartureTime,
			StopID:        record.StopID,
			StopSequence:  record.StopSequence,
		}
	}
	return out, nil
}

func (r *BusRepository) GetStopByID(ctx context.Context, stopID string) (domain.Stop, error) {
	var record model.Stop
	if err := r.db.WithContext(ctx).Where("stop_id = ?", stopID).First(&record).Error; err != nil {
		return domain.Stop{}, err
	}
	return domain.Stop{ID: record.ID.String(), StopID: record.StopID, StopName: record.StopName}, nil
}

func (r *BusRepository) listStopTimesForTrip(ctx context.Context, tripID string) ([]model.StopTime, error) {
	var stopTimes []model.StopTime
	if err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Preload("Stop").
		Order("stop_sequence ASC").
		Find(&stopTimes).Error; err != nil {
		return nil, err
	}
	return stopTimes, nil
}

func (r *BusRepository) isTripActive(ctx context.Context, serviceID string, date time.Time) (bool, error) {
	dateKey := date.Format("20060102")

	var dateException model.CalendarDate
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND date = ?", serviceID, dateKey).
		First(&dateException).Error
	if err == nil {
		return dateException.ExceptionType != 0, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	var calendar model.Calendar
	err = r.db.WithContext(ctx).
		Where("service_id = ? AND start_date <= ? AND end_date >= ?", serviceID, dateKey, dateKey).
		First(&calendar).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	day := map[time.Weekday]int{
		time.Sunday:    calendar.Sunday,
		time.Monday:    calendar.Monday,
		time.Tuesday:   calendar.Tuesday,
		time.Wednesday: calendar.Wednesday,
		time.Thursday:  calendar.Thursday,
		time.Friday:    calendar.Friday,
		time.Saturday:  calendar.Saturday,
	}[date.Weekday()]
	return day == 1, nil
}

func gtfsTimeToDate(date time.Time, raw string) (time.Time, error) {
	parsed, err := parseStopTime(raw)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, time.UTC), nil
}

func parseStopTime(raw string) (time.Time, error) {
	parsed, err := time.Parse("15:04:05", raw)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
