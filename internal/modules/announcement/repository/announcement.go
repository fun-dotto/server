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

func (r *announcementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	var dbAnnouncements []database.Announcement
	dbQuery := r.db

	if query.FilterIsActive {
		dbQuery = dbQuery.Where("is_active = ?", true)
	}

	sortDateDirection := func() string {
		switch query.SortByDate {
		case domain.SortDirectionAsc:
			return "ASC"
		case domain.SortDirectionDesc:
			return "DESC"
		default:
			return "ASC"
		}
	}()

	dbQuery = dbQuery.Order("date " + sortDateDirection)

	if err := dbQuery.Find(&dbAnnouncements).Error; err != nil {
		return nil, err
	}

	domainAnnouncements := make([]domain.Announcement, len(dbAnnouncements))
	for i, dbAnnouncement := range dbAnnouncements {
		domainAnnouncements[i] = dbAnnouncement.ToDomain()
	}

	return domainAnnouncements, nil
}
