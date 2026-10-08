package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewRoomReservationRepository(db *gorm.DB) *Repository[model.RoomReservation] {
	return newRepository[model.RoomReservation](db, []string{"id"})
}
