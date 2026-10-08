package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewRouteRepository(db *gorm.DB) *Repository[model.Route] {
	return newRepository[model.Route](db, []string{"route_id"})
}
