package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewFareRuleRepository(db *gorm.DB) *Repository[model.FareRule] {
	return newRepository[model.FareRule](db, []string{"id"})
}
