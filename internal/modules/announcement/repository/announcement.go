package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/announcement/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *announcementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	var records []model.Announcement
	dbQuery := r.db.WithContext(ctx)

	if query.FilterIsActive {
		dbQuery = dbQuery.Where("available_from <= NOW()").Where("available_until IS NULL OR available_until > NOW()")
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

	dbQuery = dbQuery.Order("available_from " + sortDateDirection)

	if err := dbQuery.Find(&records).Error; err != nil {
		return nil, err
	}

	announcements := make([]domain.Announcement, len(records))
	for i, record := range records {
		announcements[i] = toDomainAnnouncement(record)
	}

	return announcements, nil
}

func (r *announcementRepository) GetAnnouncementByID(ctx context.Context, id string) (domain.Announcement, error) {
	var record model.Announcement
	if err := r.db.WithContext(ctx).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Announcement{}, domain.ErrNotFound
		}
		return domain.Announcement{}, err
	}
	return toDomainAnnouncement(record), nil
}

func (r *announcementRepository) CreateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	record := fromDomainAnnouncement(announcement)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Announcement{}, err
	}
	return toDomainAnnouncement(record), nil
}

func (r *announcementRepository) UpdateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	record := fromDomainAnnouncement(announcement)
	result := r.db.WithContext(ctx).Model(&model.Announcement{}).Where("id = ?", parseUUIDOrNil(announcement.ID)).Updates(map[string]any{
		"title":           record.Title,
		"url":             record.URL,
		"available_from":  record.AvailableFrom,
		"available_until": record.AvailableUntil,
	})
	if result.Error != nil {
		return domain.Announcement{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.Announcement{}, domain.ErrNotFound
	}
	return r.GetAnnouncementByID(ctx, announcement.ID)
}

func (r *announcementRepository) DeleteAnnouncement(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", parseUUIDOrNil(id)).Delete(&model.Announcement{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
