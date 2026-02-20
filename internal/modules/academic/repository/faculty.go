package repository

import (
	"context"

	"github.com/fun-dotto/subject-api/internal/database"
	"github.com/fun-dotto/subject-api/internal/domain"
	"gorm.io/gorm"
)

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{db: db}
}

func (r *FacultyRepository) List(ctx context.Context) ([]domain.Faculty, error) {
	var records []database.Faculty
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.Faculty, len(records))
	for i, rec := range records {
		results[i] = database.FacultyToDomain(rec)
	}
	return results, nil
}

func (r *FacultyRepository) GetByID(ctx context.Context, id string) (domain.Faculty, error) {
	var record database.Faculty
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return domain.Faculty{}, err
	}
	return database.FacultyToDomain(record), nil
}

func (r *FacultyRepository) Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	record := database.FacultyFromDomain(faculty)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Faculty{}, err
	}
	return database.FacultyToDomain(record), nil
}

func (r *FacultyRepository) Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	record := database.FacultyFromDomain(faculty)
	if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
		return domain.Faculty{}, err
	}
	return database.FacultyToDomain(record), nil
}

func (r *FacultyRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Faculty{}, "id = ?", id).Error
}
