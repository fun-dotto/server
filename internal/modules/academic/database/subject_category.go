package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type SubjectCategory struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SubjectCategoryToDomain(m SubjectCategory) domain.SubjectCategory {
	return domain.SubjectCategory{
		ID:   m.ID,
		Name: m.Name,
	}
}

func SubjectCategoryFromDomain(d domain.SubjectCategory) SubjectCategory {
	return SubjectCategory{
		ID:   d.ID,
		Name: d.Name,
	}
}
