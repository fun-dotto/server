package repository

import (
	"context"

	"github.com/fun-dotto/subject-api/internal/database"
	"github.com/fun-dotto/subject-api/internal/domain"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) List(ctx context.Context) ([]domain.Course, error) {
	var records []database.Course
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.Course, len(records))
	for i, rec := range records {
		results[i] = database.CourseToDomain(rec)
	}
	return results, nil
}

func (r *CourseRepository) GetByID(ctx context.Context, id string) (domain.Course, error) {
	var record database.Course
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return domain.Course{}, err
	}
	return database.CourseToDomain(record), nil
}

func (r *CourseRepository) Create(ctx context.Context, course domain.Course) (domain.Course, error) {
	record := database.CourseFromDomain(course)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Course{}, err
	}
	return database.CourseToDomain(record), nil
}

func (r *CourseRepository) Update(ctx context.Context, course domain.Course) (domain.Course, error) {
	record := database.CourseFromDomain(course)
	if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
		return domain.Course{}, err
	}
	return database.CourseToDomain(record), nil
}

func (r *CourseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Course{}, "id = ?", id).Error
}
