package repository

import "gorm.io/gorm"

type FacultyRoomRepository struct {
	db *gorm.DB
}

func NewFacultyRoomRepository(db *gorm.DB) *FacultyRoomRepository {
	return &FacultyRoomRepository{db: db}
}
