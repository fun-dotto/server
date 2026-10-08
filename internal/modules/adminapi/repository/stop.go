package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewStopRepository(db *gorm.DB) *Repository[model.Stop] {
	return newRepository[model.Stop](db, []string{"stop_id"})
}
