package database

import (
	"github.com/fun-dotto/api-template/internal/domain"
)

type User struct {
	ID     string  `gorm:"type:text;primaryKey"`
	Email  string  `gorm:"type:text;not null"`
	Grade  *string `gorm:"type:text"`
	Course *string `gorm:"type:text"`
	Class  *string `gorm:"type:text"`
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		ID:     u.ID,
		Email:  u.Email,
		Grade:  u.Grade,
		Course: u.Course,
		Class:  u.Class,
	}
}

func UserFromDomain(u domain.User) User {
	return User{
		ID:     u.ID,
		Email:  u.Email,
		Grade:  u.Grade,
		Course: u.Course,
		Class:  u.Class,
	}
}
