package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) *Repository[model.User] {
	return newRepository[model.User](db, []string{"id"})
}
