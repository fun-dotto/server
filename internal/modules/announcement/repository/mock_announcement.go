package repository

import (
	"sort"
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
				Title:    "Old Announcement",
				Date:     time.Now().Add(-24 * time.Hour),
				URL:      "https://example.com/old",
				IsActive: true,
			},
			{
				ID:       "2",
				Title:    "New Announcement",
				Date:     time.Now(),
				URL:      "https://example.com/new",
				IsActive: true,
			},
			{
				ID:       "3",
				Title:    "Inactive Announcement",
				Date:     time.Now().Add(-12 * time.Hour),
				URL:      "https://example.com/inactive",
				IsActive: false,
			},
		},
	}
}

func (m *MockAnnouncementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	result := make([]domain.Announcement, len(m.announcements))
	copy(result, m.announcements)

	if query.FilterIsActive {
		var filtered []domain.Announcement
		for _, announcement := range result {
			if announcement.IsActive {
				filtered = append(filtered, announcement)
			}
		}
		result = filtered
	}

	sort.Slice(result, func(i, j int) bool {
		if query.SortByDate == domain.SortDirectionDesc {
			return result[i].Date.After(result[j].Date)
		}
		return result[i].Date.Before(result[j].Date)
	})
	return result, nil
}
