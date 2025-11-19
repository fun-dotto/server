package repository

import (
	"time"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type MockAnnouncementRepository struct {
	announcements []domain.Announcement
}

func NewMockAnnouncementRepository() *MockAnnouncementRepository {
	return &MockAnnouncementRepository{
		announcements: []domain.Announcement{
			{
				ID:       "1",
				Title:    "Announcement 1",
				Date:     time.Now(),
				URL:      "https://example.com",
				IsActive: true,
			},
		},
	}
}

func (m *MockAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
	return m.announcements, nil
}
