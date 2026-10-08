package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewTripRepository(db *gorm.DB) *Repository[model.Trip] {
	return newRepository[model.Trip](db, []string{"trip_id"})
}
