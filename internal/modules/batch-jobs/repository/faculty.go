package repository

import "gorm.io/gorm"

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{db: db}
}
