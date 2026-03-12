package repository

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) subjectPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Faculty").
		Preload("DayOfWeekTimetableSlots").
		Preload("Categories").
		Preload("EligibleAttributes").
		Preload("Requirements").
		Preload("Requirements.Course")
}

func (r *SubjectRepository) List(ctx context.Context) ([]domain.Subject, error) {
	var records []database.Subject
	if err := r.subjectPreload(r.db.WithContext(ctx)).Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.Subject, len(records))
	for i, rec := range records {
		results[i] = database.SubjectToDomain(rec)
	}
	return results, nil
}

func (r *SubjectRepository) GetByID(ctx context.Context, id string) (domain.Subject, error) {
	var record database.Subject
	if err := r.subjectPreload(r.db.WithContext(ctx)).First(&record, "id = ?", id).Error; err != nil {
		return domain.Subject{}, err
	}
	return database.SubjectToDomain(record), nil
}

func (r *SubjectRepository) Create(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	record := database.SubjectFromDomain(subject)

	for i := range record.EligibleAttributes {
		record.EligibleAttributes[i].ID = uuid.New().String()
		record.EligibleAttributes[i].SubjectID = record.ID
	}
	for i := range record.Requirements {
		record.Requirements[i].ID = uuid.New().String()
		record.Requirements[i].SubjectID = record.ID
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("DayOfWeekTimetableSlots", "Categories").Create(&record).Error; err != nil {
			return err
		}
		if len(record.DayOfWeekTimetableSlots) > 0 {
			if err := tx.Model(&record).Association("DayOfWeekTimetableSlots").Replace(record.DayOfWeekTimetableSlots); err != nil {
				return err
			}
		}
		if len(record.Categories) > 0 {
			if err := tx.Model(&record).Association("Categories").Replace(record.Categories); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Subject{}, err
	}

	return r.GetByID(ctx, record.ID)
}

func (r *SubjectRepository) Update(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	record := database.SubjectFromDomain(subject)

	for i := range record.EligibleAttributes {
		record.EligibleAttributes[i].ID = uuid.New().String()
		record.EligibleAttributes[i].SubjectID = record.ID
	}
	for i := range record.Requirements {
		record.Requirements[i].ID = uuid.New().String()
		record.Requirements[i].SubjectID = record.ID
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&record).Omit("DayOfWeekTimetableSlots", "Categories", "EligibleAttributes", "Requirements").Save(&record).Error; err != nil {
			return err
		}

		// Replace M:N associations
		if err := tx.Model(&record).Association("DayOfWeekTimetableSlots").Replace(record.DayOfWeekTimetableSlots); err != nil {
			return err
		}
		if err := tx.Model(&record).Association("Categories").Replace(record.Categories); err != nil {
			return err
		}

		// Replace 1:N owned records: delete old, create new
		if err := tx.Where("subject_id = ?", record.ID).Delete(&database.SubjectEligibleAttribute{}).Error; err != nil {
			return err
		}
		if len(record.EligibleAttributes) > 0 {
			if err := tx.Create(&record.EligibleAttributes).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("subject_id = ?", record.ID).Delete(&database.SubjectRequirement{}).Error; err != nil {
			return err
		}
		if len(record.Requirements) > 0 {
			if err := tx.Create(&record.Requirements).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Subject{}, err
	}

	return r.GetByID(ctx, record.ID)
}

func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record database.Subject
		if err := tx.First(&record, "id = ?", id).Error; err != nil {
			return err
		}

		// Clear M:N associations
		if err := tx.Model(&record).Association("DayOfWeekTimetableSlots").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&record).Association("Categories").Clear(); err != nil {
			return err
		}

		// Delete 1:N owned records
		if err := tx.Where("subject_id = ?", id).Delete(&database.SubjectEligibleAttribute{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subject_id = ?", id).Delete(&database.SubjectRequirement{}).Error; err != nil {
			return err
		}

		return tx.Delete(&record).Error
	})
}
