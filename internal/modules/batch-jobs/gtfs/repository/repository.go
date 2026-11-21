package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/gtfs/domain"
	"gorm.io/gorm"
)

func Save(ctx context.Context, db *gorm.DB, schedule domain.Schedule) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := clearGTFSData(tx); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "calendars", schedule.Calendars); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "zones", schedule.Zones); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "fare_attributes", schedule.FareAttributes); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "routes", schedule.Routes); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "stops", schedule.Stops); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "trips", schedule.Trips); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "calendar_dates", schedule.CalendarDates); err != nil {
			return err
		}
		if err := createGTFSBatch(tx, "fare_rules", schedule.FareRules); err != nil {
			return err
		}
		return createGTFSBatch(tx, "stop_times", schedule.StopTimes)
	})
}

func clearGTFSData(tx *gorm.DB) error {
	tables := []string{
		"stop_times",
		"calendar_dates",
		"fare_rules",
		"trips",
		"stops",
		"fare_attributes",
		"routes",
		"zones",
		"calendars",
	}
	for _, table := range tables {
		if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
			return fmt.Errorf("delete %s: %w", table, err)
		}
	}
	fmt.Println("Deleted existing GTFS data")
	return nil
}

func createGTFSBatch[T any](tx *gorm.DB, table string, rows []T) error {
	if len(rows) == 0 {
		return nil
	}
	if err := tx.CreateInBatches(&rows, 1000).Error; err != nil {
		return fmt.Errorf("create %s: %w", table, err)
	}
	return nil
}
