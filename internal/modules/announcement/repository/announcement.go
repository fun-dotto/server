package repository

import (
	"github.com/fun-dotto/announcement-api/internal/database"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"gorm.io/gorm"
)

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *announcementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) GetAnnouncements(isActive *bool) ([]domain.Announcement, error) {
	var dbAnnouncements []database.Announcement
	query := r.db

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Find(&dbAnnouncements).Error; err != nil {
		return nil, err
	}

	domainAnnouncements := make([]domain.Announcement, len(dbAnnouncements))
	for i, dbAnnouncement := range dbAnnouncements {
		domainAnnouncements[i] = dbAnnouncement.ToDomain()
	}

	return domainAnnouncements, nil
}
