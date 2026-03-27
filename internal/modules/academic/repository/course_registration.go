package repository

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRegistrationRepository struct {
	db *gorm.DB
}

func NewCourseRegistrationRepository(db *gorm.DB) *CourseRegistrationRepository {
	return &CourseRegistrationRepository{db: db}
}

func (r *CourseRegistrationRepository) courseRegistrationPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subject.Faculties").
		Preload("Subject.EligibleAttributes").
		Preload("Subject.Requirements")
}

func (r *CourseRegistrationRepository) List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error) {
	query := r.courseRegistrationPreload(r.db.WithContext(ctx)).
		Joins("JOIN subjects ON subjects.id = course_registrations.subject_id").
		Where("course_registrations.user_id = ?", filter.UserID)

	if filter.Year != nil {
		query = query.Where("subjects.year = ?", *filter.Year)
	}
	if len(filter.Semesters) > 0 {
		semesters := make([]string, len(filter.Semesters))
		for i, s := range filter.Semesters {
			semesters[i] = string(s)
		}
		query = query.Where("subjects.semester IN ?", semesters)
	}

	var records []database.CourseRegistration
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	items := make([]domain.CourseRegistration, len(records))
	for i, rec := range records {
		items[i] = database.CourseRegistrationToDomain(rec)
	}
	return items, nil
}

func (r *CourseRegistrationRepository) Create(ctx context.Context, cr domain.CourseRegistration) (domain.CourseRegistration, error) {
	record := database.CourseRegistrationFromDomain(cr)
	record.ID = uuid.New().String()

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.CourseRegistration{}, err
	}

	var created database.CourseRegistration
	if err := r.courseRegistrationPreload(r.db.WithContext(ctx)).First(&created, "id = ?", record.ID).Error; err != nil {
		return domain.CourseRegistration{}, err
	}
	return database.CourseRegistrationToDomain(created), nil
}

func (r *CourseRegistrationRepository) Delete(ctx context.Context, id string) error {
	var record database.CourseRegistration
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Delete(&record)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
