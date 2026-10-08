package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewFacultyRepository(db *gorm.DB) *Repository[model.Faculty] {
	return newRepository[model.Faculty](db, []string{"id"})
}
