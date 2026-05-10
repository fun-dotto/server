package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type CancelledClassRepository struct {
	db *gorm.DB
}

func NewCancelledClassRepository(db *gorm.DB) *CancelledClassRepository {
	return &CancelledClassRepository{db: db}
}

func (r *CancelledClassRepository) cancelledClassPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subject.Faculties.Faculty").
		Preload("Subject.EligibleAttributes").
		Preload("Subject.Requirements")
}

func (r *CancelledClassRepository) List(ctx context.Context, filter domain.CancelledClassListFilter) ([]domain.CancelledClass, error) {
	var records []model.CancelledClass
	query := r.cancelledClassPreload(r.db.WithContext(ctx))

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
		Joins("JOIN subjects ON subjects.id = cancelled_classes.subject_id").
		Order("cancelled_classes.date ASC, cancelled_classes.period ASC, subjects.syllabus_id ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.CancelledClass, len(records))
	for i, rec := range records {
		results[i] = cancelledClassToDomain(rec)
	}
	return results, nil
}

func (r *CancelledClassRepository) GetByID(ctx context.Context, id string) (domain.CancelledClass, error) {
	var record model.CancelledClass
	if err := r.cancelledClassPreload(r.db.WithContext(ctx)).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.CancelledClass{}, err
	}
	return cancelledClassToDomain(record), nil
}

func (r *CancelledClassRepository) Create(ctx context.Context, cc domain.CancelledClass) (domain.CancelledClass, error) {
	record := cancelledClassFromDomain(cc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.CancelledClass{}, err
	}
	return r.GetByID(ctx, record.ID.String())
}

func (r *CancelledClassRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.CancelledClass
		if err := tx.Where("id = ?", uid).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		result := tx.Where("id = ?", uid).Delete(&model.CancelledClass{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
