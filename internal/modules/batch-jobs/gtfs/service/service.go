package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/gtfs/fetcher"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/gtfs/parser"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/gtfs/repository"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Import(ctx context.Context, url string) error {
	tempDir, err := os.MkdirTemp("", "gtfs-import-")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	extractDir := filepath.Join(tempDir, "extracted")
	if err := fetcher.DownloadAndExtract(url, extractDir); err != nil {
		return err
	}

	schedule, err := parser.Parse(extractDir)
	if err != nil {
		return fmt.Errorf("convert GTFS files to models: %w", err)
	}
	if err := repository.Save(ctx, s.db, schedule); err != nil {
		return fmt.Errorf("store GTFS models: %w", err)
	}

	fmt.Println("stored GTFS models:")
	fmt.Printf(
		"calendars=%d\ncalendar_dates=%d\nzones=%d\nfare_attributes=%d\nfare_rules=%d\nroutes=%d\nstops=%d\ntrips=%d\nstop_times=%d\n",
		len(schedule.Calendars),
		len(schedule.CalendarDates),
		len(schedule.Zones),
		len(schedule.FareAttributes),
		len(schedule.FareRules),
		len(schedule.Routes),
		len(schedule.Stops),
		len(schedule.Trips),
		len(schedule.StopTimes),
	)
	return nil
}
