package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type MakeupClassRepository struct {
	db *gorm.DB
}

func NewMakeupClassRepository(db *gorm.DB) *MakeupClassRepository {
	return &MakeupClassRepository{db: db}
}

func (r *MakeupClassRepository) makeupClassPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subject.Faculties.Faculty").
		Preload("Subject.EligibleAttributes").
		Preload("Subject.Requirements")
}

func (r *MakeupClassRepository) List(ctx context.Context, filter domain.MakeupClassListFilter) ([]domain.MakeupClass, error) {
	var records []model.MakeupClass
	query := r.makeupClassPreload(r.db.WithContext(ctx))

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
		Joins("JOIN subjects ON subjects.id = makeup_classes.subject_id").
		Order("makeup_classes.date ASC, makeup_classes.period ASC, subjects.syllabus_id ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.MakeupClass, len(records))
	for i, rec := range records {
		results[i] = makeupClassToDomain(rec)
	}
	return results, nil
}

func (r *MakeupClassRepository) GetByID(ctx context.Context, id string) (domain.MakeupClass, error) {
	var record model.MakeupClass
	if err := r.makeupClassPreload(r.db.WithContext(ctx)).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.MakeupClass{}, err
	}
	return makeupClassToDomain(record), nil
}

func (r *MakeupClassRepository) Create(ctx context.Context, mc domain.MakeupClass) (domain.MakeupClass, error) {
	record := makeupClassFromDomain(mc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.MakeupClass{}, err
	}
	return r.GetByID(ctx, record.ID.String())
}

func (r *MakeupClassRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.MakeupClass
		if err := tx.Where("id = ?", uid).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		result := tx.Where("id = ?", uid).Delete(&model.MakeupClass{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
