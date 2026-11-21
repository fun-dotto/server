package repository

import (
	"context"
	"errors"
	"uuid"

	"gorm.io/gorm"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *NotificationRepository) DeleteNotification(ctx context.Context, id string) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrNotFound
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Notification
		if err := tx.First(&existing, "id = ?", notificationID.String()).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		if err := tx.Where("notification_id = ?", notificationID.String()).Delete(&model.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		return tx.Delete(&existing).Error
	})
}
