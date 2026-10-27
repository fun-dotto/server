package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewFacultyRoomRepository(db *gorm.DB) *Repository[model.FacultyRoom] {
	return newRepository[model.FacultyRoom](db, []string{"faculty_id", "room_id", "year"})
}
