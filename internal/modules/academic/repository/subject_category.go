package repository

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"gorm.io/gorm"
)

type SubjectCategoryRepository struct {
	db *gorm.DB
}

func NewSubjectCategoryRepository(db *gorm.DB) *SubjectCategoryRepository {
	return &SubjectCategoryRepository{db: db}
}

func (r *SubjectCategoryRepository) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	var records []database.SubjectCategory
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.SubjectCategory, len(records))
	for i, rec := range records {
		results[i] = database.SubjectCategoryToDomain(rec)
	}
	return results, nil
}

func (r *SubjectCategoryRepository) GetByID(ctx context.Context, id string) (domain.SubjectCategory, error) {
	var record database.SubjectCategory
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return domain.SubjectCategory{}, err
	}
	return database.SubjectCategoryToDomain(record), nil
}

func (r *SubjectCategoryRepository) Create(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error) {
	record := database.SubjectCategoryFromDomain(category)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.SubjectCategory{}, err
	}
	return database.SubjectCategoryToDomain(record), nil
}

func (r *SubjectCategoryRepository) Update(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error) {
	record := database.SubjectCategoryFromDomain(category)
	if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
		return domain.SubjectCategory{}, err
	}
	return database.SubjectCategoryToDomain(record), nil
}

func (r *SubjectCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.SubjectCategory{}, "id = ?", id).Error
}
