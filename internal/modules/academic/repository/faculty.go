package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{db: db}
}

func (r *FacultyRepository) List(ctx context.Context, ids []string) ([]domain.Faculty, error) {
	var records []model.Faculty
	query := r.db.WithContext(ctx)
	if len(ids) > 0 {
		query = query.Where("id IN ?", parseUUIDs(ids))
	}
	if err := query.Order("email ASC").Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.Faculty, len(records))
	for i, rec := range records {
		results[i] = facultyToDomain(rec)
	}
	return results, nil
}

func (r *FacultyRepository) GetByID(ctx context.Context, id string) (domain.Faculty, error) {
	var record model.Faculty
	if err := r.db.WithContext(ctx).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.Faculty{}, err
	}
	return facultyToDomain(record), nil
}

func (r *FacultyRepository) Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	record := facultyFromDomain(faculty)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Faculty{}, err
	}
	return facultyToDomain(record), nil
}

func (r *FacultyRepository) Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	id := parseUUIDOrNil(faculty.ID)
	if err := r.db.WithContext(ctx).Model(&model.Faculty{}).Where("id = ?", id).Updates(map[string]any{
		"name":  faculty.Name,
		"email": faculty.Email,
	}).Error; err != nil {
		return domain.Faculty{}, err
	}
	return r.GetByID(ctx, faculty.ID)
}

func (r *FacultyRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var faculty model.Faculty
		if err := tx.Where("id = ?", uid).First(&faculty).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		if err := tx.Where("faculty_id = ?", uid).Delete(&model.SubjectFaculty{}).Error; err != nil {
			return err
		}

		if err := tx.Where("faculty_id = ?", uid).Delete(&model.FacultyRoom{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", uid).Delete(&model.Faculty{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
