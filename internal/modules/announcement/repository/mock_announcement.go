package repository

import (
	"context"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type MockAnnouncementRepository struct {
	GetAnnouncementsFunc func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error)
}

func (m *MockAnnouncementRepository) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	if m.GetAnnouncementsFunc != nil {
		return m.GetAnnouncementsFunc(ctx, query)
	}
	return []domain.Announcement{}, nil
}
