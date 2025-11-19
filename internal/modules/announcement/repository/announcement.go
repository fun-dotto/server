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

func (r *announcementRepository) GetAnnouncements() ([]domain.Announcement, error) {
	var announcements []database.AnnouncementModel
	if err := r.db.Find(&announcements).Error; err != nil {
		return nil, err
	}

	announcementDomains := make([]domain.Announcement, len(announcements))
	for i, announcement := range announcements {
		announcementDomains[i] = announcement.ToDomain()
	}

	return announcementDomains, nil
}
