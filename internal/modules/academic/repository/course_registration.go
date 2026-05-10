package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
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
		Preload("Subject.Faculties.Faculty").
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

	var records []model.CourseRegistration
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.CourseRegistration, len(records))
	for i, rec := range records {
		results[i] = courseRegistrationToDomain(rec)
	}
	return results, nil
}

func (r *CourseRegistrationRepository) Create(ctx context.Context, cr domain.CourseRegistration) (domain.CourseRegistration, error) {
	record := courseRegistrationFromDomain(cr)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.CourseRegistration{}, err
	}

	var created model.CourseRegistration
	if err := r.courseRegistrationPreload(r.db.WithContext(ctx)).
		Where("user_id = ? AND subject_id = ?", record.UserID, record.SubjectID).
		First(&created).Error; err != nil {
		return domain.CourseRegistration{}, err
	}
	return courseRegistrationToDomain(created), nil
}

func (r *CourseRegistrationRepository) Delete(ctx context.Context, id string) error {
	userID, subjectID, err := decodeCourseRegistrationID(id)
	if err != nil {
		return gorm.ErrRecordNotFound
	}

	var record model.CourseRegistration
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND subject_id = ?", userID, subjectID).
		First(&record).Error; err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND subject_id = ?", userID, subjectID).
		Delete(&model.CourseRegistration{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
