package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type Faculty struct {
	ID    string `gorm:"primaryKey;type:uuid"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func FacultyToDomain(m Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
}

func FacultyFromDomain(faculty domain.Faculty) Faculty {
	return Faculty{
		ID:    faculty.ID,
		Name:  faculty.Name,
		Email: faculty.Email,
	}
}
