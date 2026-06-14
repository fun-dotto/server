package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	notification.ID = uuid.New().String()
	notificationID, err := uuid.Parse(notification.ID)
	if err != nil {
		return domain.Notification{}, err
	}

	dbNotification := notificationFromDomain(notification)

	uniqueTargets := uniqueTargetUsers(notification.TargetUsers)

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbNotification).Error; err != nil {
			return err
		}

		if len(uniqueTargets) > 0 {
			targets := make([]model.NotificationTargetUser, 0, len(uniqueTargets))
			for _, t := range uniqueTargets {
				targets = append(targets, model.NotificationTargetUser{
					NotificationID: notificationID,
					UserID:         t.UserID,
					NotifiedAt:     t.NotifiedAt,
				})
			}
			// Notification/User の関連は親側で作成済みなので Omit して GORM の auto-upsert を抑止する。
			if err := tx.Omit("Notification", "User").Create(&targets).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Notification{}, err
	}

	return notificationToDomain(dbNotification, uniqueTargets), nil
}
