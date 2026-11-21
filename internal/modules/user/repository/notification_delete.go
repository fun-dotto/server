package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *NotificationRepository) DeleteNotification(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Notification
		if err := tx.First(&existing, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		if err := tx.Where("notification_id = ?", id).Delete(&model.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		return tx.Delete(&existing).Error
	})
}
