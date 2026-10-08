package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewFareAttributeRepository(db *gorm.DB) *Repository[model.FareAttribute] {
	return newRepository[model.FareAttribute](db, []string{"fare_id"})
}
