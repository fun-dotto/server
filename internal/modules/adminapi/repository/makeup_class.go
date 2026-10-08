package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewMakeupClassRepository(db *gorm.DB) *Repository[model.MakeupClass] {
	return newRepository[model.MakeupClass](db, []string{"id"})
}
