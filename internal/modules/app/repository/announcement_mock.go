package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// MockAnnouncementRepository はテスト用のモック
type MockAnnouncementRepository struct {
	announcements       []domain.Announcement
	getAnnouncementsErr error
}

func NewMockAnnouncementRepository() *MockAnnouncementRepository {
	return &MockAnnouncementRepository{
		announcements: []domain.Announcement{
			{
				ID:            "1",
				Title:         "Announcement 1",
				AvailableFrom: time.Now(),
				URL:           "https://example.com",
			},
		},
	}
}

func NewMockAnnouncementRepositoryWithError(err error) *MockAnnouncementRepository {
	m := NewMockAnnouncementRepository()
	m.getAnnouncementsErr = err
	return m
}

func (m *MockAnnouncementRepository) GetAnnouncements(_ context.Context, _ domain.AnnouncementQuery) ([]domain.Announcement, error) {
	if m.getAnnouncementsErr != nil {
		return nil, m.getAnnouncementsErr
	}
	return m.announcements, nil
}
