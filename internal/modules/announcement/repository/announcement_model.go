package repository

import (
	"time"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type AnnouncementModel struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Title     string    `gorm:"type:varchar(500);not null"`
	Date      time.Time `gorm:"not null;index"`
	URL       string    `gorm:"type:varchar(1000);not null"`
	IsActive  bool      `gorm:"not null;default:true;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (AnnouncementModel) TableName() string {
	return "announcements"
}

func (m *AnnouncementModel) ToDomain() domain.Announcement {
	return domain.Announcement{
		ID:       m.ID,
		Title:    m.Title,
		Date:     m.Date,
		URL:      m.URL,
		IsActive: m.IsActive,
	}
}

func FromDomain(announcement domain.Announcement) AnnouncementModel {
	return AnnouncementModel{
		ID:       announcement.ID,
		Title:    announcement.Title,
		Date:     announcement.Date,
		URL:      announcement.URL,
		IsActive: announcement.IsActive,
	}
}
