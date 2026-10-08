package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewNotificationRepository(db *gorm.DB) *Repository[model.Notification] {
	return newRepository[model.Notification](db, []string{"id"})
}
