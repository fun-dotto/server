package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) DeleteNotification(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing database.Notification
		if err := tx.First(&existing, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		if err := tx.Where("notification_id = ?", id).Delete(&database.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		return tx.Delete(&existing).Error
	})
}
