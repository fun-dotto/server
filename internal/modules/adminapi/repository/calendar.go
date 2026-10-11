package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewCalendarRepository(db *gorm.DB) *Repository[model.Calendar] {
	return newRepository[model.Calendar](db, []string{"service_id"})
}
