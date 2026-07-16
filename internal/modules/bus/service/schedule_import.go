package service

import (
	"context"
	"fmt"
)

// Import は URL から GTFS を取得・整形し、DB を全置換する。
func (s *ScheduleService) Import(ctx context.Context, url string) error {
	archive, err := s.FetchArchive(ctx, url)
	if err != nil {
		return err
	}

	stops, routes, calendars, calendarDates, trips, stopTimes, fareRules, err := s.ConvertToDomain(archive)
	if err != nil {
		return fmt.Errorf("convert gtfs to domain: %w", err)
	}

	if err := s.repo.ReplaceAll(ctx, stops, routes, calendars, calendarDates, trips, stopTimes, fareRules); err != nil {
		return fmt.Errorf("replace gtfs: %w", err)
	}
	return nil
}
