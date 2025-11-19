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
				Title:    "Active Announcement",
				Date:     time.Now(),
				URL:      "https://example.com/active",
				IsActive: true,
			},
			{
				ID:       "2",
				Title:    "Inactive Announcement",
				Date:     time.Now(),
				URL:      "https://example.com/inactive",
				IsActive: false,
			},
		},
	}
}

func (m *MockAnnouncementRepository) GetAnnouncements(isActive *bool) ([]domain.Announcement, error) {
	if isActive == nil {
		return m.announcements, nil
	}
	
	var filtered []domain.Announcement
	for _, announcement := range m.announcements {
		if announcement.IsActive == *isActive {
			filtered = append(filtered, announcement)
		}
	}
	return filtered, nil
}
