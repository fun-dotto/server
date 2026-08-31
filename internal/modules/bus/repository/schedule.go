package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// insertBatchSize は stop_times / fare_rules のような大量行を分割 INSERT する単位。
const insertBatchSize = 1000

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// ReplaceAll は既存の GTFS データを全削除し、渡された世代で置き換える。
// 1 トランザクション内で行い、途中失敗時は前の世代を保つ。
// 外部キー制約に合わせ、削除は子→親、挿入は親→子の順で行う。
func (r *ScheduleRepository) ReplaceAll(
	ctx context.Context,
	stops []domain.Stop,
	routes []domain.Route,
	calendars []domain.Calendar,
	calendarDates []domain.CalendarDate,
	trips []domain.Trip,
	stopTimes []domain.StopTime,
	fareRules []domain.FareRule,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := deleteAll(tx); err != nil {
			return err
		}
		return insertAll(tx, stops, routes, calendars, calendarDates, trips, stopTimes, fareRules)
	})
}

func deleteAll(tx *gorm.DB) error {
	// 子 → 親 の順。
	models := []any{
		&model.StopTime{},
		&model.FareRule{},
		&model.Trip{},
		&model.CalendarDate{},
		&model.Stop{},
		&model.Route{},
		&model.Calendar{},
	}
	for _, m := range models {
		if err := tx.Where("1 = 1").Delete(m).Error; err != nil {
			return fmt.Errorf("delete %T: %w", m, err)
		}
	}
	return nil
}

func insertAll(
	tx *gorm.DB,
	stops []domain.Stop,
	routes []domain.Route,
	calendars []domain.Calendar,
	calendarDates []domain.CalendarDate,
	trips []domain.Trip,
	stopTimes []domain.StopTime,
	fareRules []domain.FareRule,
) error {
	// 親（他を参照しない）。
	if err := insertBatch(tx, stopsFromDomain(stops), "stops"); err != nil {
		return err
	}
	if err := insertBatch(tx, routesFromDomain(routes), "routes"); err != nil {
		return err
	}
	if err := insertBatch(tx, calendarsFromDomain(calendars), "calendars"); err != nil {
		return err
	}

	// 親を参照する中間。
	if err := insertBatch(tx, calendarDatesFromDomain(calendarDates), "calendar_dates"); err != nil {
		return err
	}
	if err := insertBatch(tx, tripsFromDomain(trips), "trips"); err != nil {
		return err
	}
	if err := insertBatch(tx, fareRulesFromDomain(fareRules), "fare_rules"); err != nil {
		return err
	}

	// trips / stops を参照する末端。
	return insertBatch(tx, stopTimesFromDomain(stopTimes), "stop_times")
}

func insertBatch[T any](tx *gorm.DB, records []T, name string) error {
	if len(records) == 0 {
		return nil
	}
	if err := tx.CreateInBatches(records, insertBatchSize).Error; err != nil {
		return fmt.Errorf("insert %s: %w", name, err)
	}
	return nil
}
