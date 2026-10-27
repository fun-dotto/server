package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewRoomChangeRepository(db *gorm.DB) *Repository[model.RoomChange] {
	return newRepository[model.RoomChange](db, []string{"id"})
}
