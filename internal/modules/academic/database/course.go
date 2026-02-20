package database

import (
	"time"

	"github.com/fun-dotto/subject-api/internal/domain"
)

type Course struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CourseToDomain(m Course) domain.Course {
	return domain.Course{
		ID:   m.ID,
		Name: m.Name,
	}
}

func CourseFromDomain(d domain.Course) Course {
	return Course{
		ID:   d.ID,
		Name: d.Name,
	}
}
