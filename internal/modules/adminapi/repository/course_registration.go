package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewCourseRegistrationRepository(db *gorm.DB) *Repository[model.CourseRegistration] {
	return newRepository[model.CourseRegistration](db, []string{"user_id", "subject_id"})
}
