package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewZoneRepository(db *gorm.DB) *Repository[model.Zone] {
	return newRepository[model.Zone](db, []string{"zone_id"})
}
