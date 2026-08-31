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

// ListTripsByDate は指定された日付の有効なバスの運行情報を取得します。
func (r *BusRepository) ListTripsByDate(ctx context.Context, date time.Time) ([]domain.Trip, error) {
	var trips []model.Trip
	if err := r.db.WithContext(ctx).Find(&trips).Error; err != nil {
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

// GetRoutesByIDs は指定された路線IDの路線情報を取得します。
func (r *BusRepository) GetRoutesByIDs(ctx context.Context, routeIDs []string) (map[string]domain.Route, error) {
	if len(routeIDs) == 0 {
		return map[string]domain.Route{}, nil
	}

	var records []model.Route
	if err := r.db.WithContext(ctx).Where("route_id IN ?", routeIDs).Find(&records).Error; err != nil {
		return nil, err
	}

	out := make(map[string]domain.Route, len(records))
	for _, record := range records {
		out[record.RouteID] = domain.Route{
			ID:             record.ID.String(),
			RouteID:        record.RouteID,
			RouteShortName: record.RouteShortName,
		}
	}
	return out, nil
}

// ListStopTimesByTripIDs は指定されたバスIDの停車時刻情報を取得します。
func (r *BusRepository) ListStopTimesByTripIDs(ctx context.Context, tripIDs []string) ([]domain.StopTime, error) {
	if len(tripIDs) == 0 {
		return []domain.StopTime{}, nil
	}

	var records []model.StopTime
	if err := r.db.WithContext(ctx).
		Where("trip_id IN ?", tripIDs).
		Order("trip_id ASC, stop_sequence ASC").
		Find(&records).Error; err != nil {
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

// GetStopsByIDs は指定された停車所IDの停車所情報を取得します。
func (r *BusRepository) GetStopsByIDs(ctx context.Context, stopIDs []string) (map[string]domain.Stop, error) {
	if len(stopIDs) == 0 {
		return map[string]domain.Stop{}, nil
	}

	var records []model.Stop
	if err := r.db.WithContext(ctx).Where("stop_id IN ?", stopIDs).Find(&records).Error; err != nil {
		return nil, err
	}

	out := make(map[string]domain.Stop, len(records))
	for _, record := range records {
		out[record.StopID] = domain.Stop{
			ID:       record.ID.String(),
			StopID:   record.StopID,
			StopName: record.StopName,
		}
	}
	return out, nil
}

// isTripActive は指定されたサービスIDと日付のバスが有効かどうかを判断します。
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
