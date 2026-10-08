package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewRoomRepository(db *gorm.DB) *Repository[model.Room] {
	return newRepository[model.Room](db, []string{"id"})
}
