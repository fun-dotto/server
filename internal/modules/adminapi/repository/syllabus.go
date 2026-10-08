package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewSyllabusRepository(db *gorm.DB) *Repository[model.Syllabus] {
	return newRepository[model.Syllabus](db, []string{"id"})
}
