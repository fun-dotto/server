package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewCancelledClassRepository(db *gorm.DB) *Repository[model.CancelledClass] {
	return newRepository[model.CancelledClass](db, []string{"id"})
}
