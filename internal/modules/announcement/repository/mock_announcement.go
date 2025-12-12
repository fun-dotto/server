package repository

import (
	"github.com/fun-dotto/announcement-api/internal/domain"
)

type MockAnnouncementRepository struct {
	GetAnnouncementsFunc func(query domain.AnnouncementQuery) ([]domain.Announcement, error)
}

func (m *MockAnnouncementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	if m.GetAnnouncementsFunc != nil {
		return m.GetAnnouncementsFunc(query)
	}
	return []domain.Announcement{}, nil
}
