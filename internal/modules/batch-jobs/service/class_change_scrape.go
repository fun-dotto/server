package service

import (
	"context"
	"time"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

// ScraperClient は学生ポータルへログインし休講/補講/部屋変更一覧の HTML を取得する。
// 実装は internal/modules/batch-jobs/scraper.Client。
type ScraperClient interface {
	FetchClassChangeHTML(ctx context.Context, userID, password string) (string, error)
}

type ScrapeSubjectRepository interface {
	ListSyllabusNames(ctx context.Context) ([]domain.SubjectRef, error)
}

type ScrapeRoomRepository interface {
	ListRooms(ctx context.Context) ([]domain.RoomRef, error)
}

type ScrapeCancelledClassRepository interface {
	Exists(ctx context.Context, subjectID string, date time.Time, period string) (bool, error)
	Insert(ctx context.Context, cc academicdomain.CancelledClass) (academicdomain.CancelledClass, error)
}

type ScrapeMakeupClassRepository interface {
	Exists(ctx context.Context, subjectID string, date time.Time, period string) (bool, error)
	Insert(ctx context.Context, mc academicdomain.MakeupClass) (academicdomain.MakeupClass, error)
}

type ScrapeRoomChangeRepository interface {
	Exists(ctx context.Context, subjectID string, date time.Time, period, originalRoomID, newRoomID string) (bool, error)
	Insert(ctx context.Context, rc academicdomain.RoomChange) (academicdomain.RoomChange, error)
}

// ClassChangeScrapeService は学生ポータルから休講/補講/部屋変更をスクレイプし、
// subjects/rooms と突合したうえで冪等に DB へ反映する。
type ClassChangeScrapeService struct {
	client     ScraperClient
	subjects   ScrapeSubjectRepository
	rooms      ScrapeRoomRepository
	cancelled  ScrapeCancelledClassRepository
	makeup     ScrapeMakeupClassRepository
	roomChange ScrapeRoomChangeRepository
}

func NewClassChangeScrapeService(
	client ScraperClient,
	subjects ScrapeSubjectRepository,
	rooms ScrapeRoomRepository,
	cancelled ScrapeCancelledClassRepository,
	makeup ScrapeMakeupClassRepository,
	roomChange ScrapeRoomChangeRepository,
) *ClassChangeScrapeService {
	return &ClassChangeScrapeService{
		client:     client,
		subjects:   subjects,
		rooms:      rooms,
		cancelled:  cancelled,
		makeup:     makeup,
		roomChange: roomChange,
	}
}
