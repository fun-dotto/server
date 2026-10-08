package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewNotificationTargetUserRepository(db *gorm.DB) *Repository[model.NotificationTargetUser] {
	return newRepository[model.NotificationTargetUser](db, []string{"notification_id", "user_id"})
}
