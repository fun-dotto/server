package repository

import (
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
	var announcements []domain.Announcement
	if err := r.db.Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}
