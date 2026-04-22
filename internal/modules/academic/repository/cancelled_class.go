package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
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
	var records []database.CancelledClass
	query := r.cancelledClassPreload(r.db.WithContext(ctx))

	if len(filter.SubjectIDs) > 0 {
		query = query.Where("subject_id IN ?", filter.SubjectIDs)
	}
	if filter.From != nil {
		query = query.Where("date >= ?", filter.From.Format("2006-01-02"))
	}
	if filter.Until != nil {
		query = query.Where("date <= ?", filter.Until.Format("2006-01-02"))
	}

	query = query.
		Joins("JOIN subjects ON subjects.id = cancelled_classes.subject_id").
		Order("cancelled_classes.date ASC, cancelled_classes.period ASC, subjects.syllabus_id ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.CancelledClass, len(records))
	for i, rec := range records {
		results[i] = database.CancelledClassToDomain(rec)
	}
	return results, nil
}

func (r *CancelledClassRepository) GetByID(ctx context.Context, id string) (domain.CancelledClass, error) {
	var record database.CancelledClass
	if err := r.cancelledClassPreload(r.db.WithContext(ctx)).First(&record, "id = ?", id).Error; err != nil {
		return domain.CancelledClass{}, err
	}
	return database.CancelledClassToDomain(record), nil
}

func (r *CancelledClassRepository) Create(ctx context.Context, cc domain.CancelledClass) (domain.CancelledClass, error) {
	dbRecord := database.CancelledClassFromDomain(cc)
	if dbRecord.ID == "" {
		dbRecord.ID = uuid.New().String()
	}

	if err := r.db.WithContext(ctx).Create(&dbRecord).Error; err != nil {
		return domain.CancelledClass{}, err
	}
	return r.GetByID(ctx, dbRecord.ID)
}

func (r *CancelledClassRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record database.CancelledClass
		if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		result := tx.Where("id = ?", id).Delete(&database.CancelledClass{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
