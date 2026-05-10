package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type RoomChangeRepository struct {
	db *gorm.DB
}

func NewRoomChangeRepository(db *gorm.DB) *RoomChangeRepository {
	return &RoomChangeRepository{db: db}
}

func (r *RoomChangeRepository) roomChangePreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subject.Faculties.Faculty").
		Preload("Subject.EligibleAttributes").
		Preload("Subject.Requirements").
		Preload("OriginalRoom").
		Preload("NewRoom")
}

func (r *RoomChangeRepository) List(ctx context.Context, filter domain.RoomChangeListFilter) ([]domain.RoomChange, error) {
	var records []model.RoomChange
	query := r.roomChangePreload(r.db.WithContext(ctx))

	if len(filter.SubjectIDs) > 0 {
		query = query.Where("subject_id IN ?", parseUUIDs(filter.SubjectIDs))
	}
	if filter.From != nil {
		query = query.Where("date >= ?", filter.From.Format(dateLayout))
	}
	if filter.Until != nil {
		query = query.Where("date <= ?", filter.Until.Format(dateLayout))
	}

	query = query.
		Joins("JOIN subjects ON subjects.id = room_changes.subject_id").
		Order("room_changes.date ASC, room_changes.period ASC, subjects.syllabus_id ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.RoomChange, len(records))
	for i, rec := range records {
		results[i] = roomChangeToDomain(rec)
	}
	return results, nil
}

func (r *RoomChangeRepository) GetByID(ctx context.Context, id string) (domain.RoomChange, error) {
	var record model.RoomChange
	if err := r.roomChangePreload(r.db.WithContext(ctx)).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.RoomChange{}, err
	}
	return roomChangeToDomain(record), nil
}

func (r *RoomChangeRepository) Create(ctx context.Context, rc domain.RoomChange) (domain.RoomChange, error) {
	record := roomChangeFromDomain(rc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.RoomChange{}, err
	}
	return r.GetByID(ctx, record.ID.String())
}

func (r *RoomChangeRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.RoomChange
		if err := tx.Where("id = ?", uid).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		result := tx.Where("id = ?", uid).Delete(&model.RoomChange{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
