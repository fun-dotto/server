package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewFCMTokenRepository(db *gorm.DB) *Repository[model.FCMToken] {
	return newRepository[model.FCMToken](db, []string{"token"})
}
