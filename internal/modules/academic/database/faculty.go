package database

import (
	"time"

	"github.com/fun-dotto/subject-api/internal/domain"
)

type Faculty struct {
	ID        string    `gorm:"type:text;primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	Email     string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FacultyToDomain(m Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
}

func FacultyFromDomain(d domain.Faculty) Faculty {
	return Faculty{
		ID:    d.ID,
		Name:  d.Name,
		Email: d.Email,
	}
}
