package service

import (
	"context"
	"net/http"
	"time"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

// scheduleRepository は取得・変換済みの GTFS を DB へ全置換する永続化境界。
type scheduleRepository interface {
	ReplaceAll(
		ctx context.Context,
		stops []domain.Stop,
		routes []domain.Route,
		calendars []domain.Calendar,
		calendarDates []domain.CalendarDate,
		trips []domain.Trip,
		stopTimes []domain.StopTime,
		fareRules []domain.FareRule,
	) error
}

type ScheduleService struct {
	repo       scheduleRepository
	httpClient *http.Client
}

func NewScheduleService(repo scheduleRepository, httpClient *http.Client) *ScheduleService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}
	return &ScheduleService{repo: repo, httpClient: httpClient}
}
