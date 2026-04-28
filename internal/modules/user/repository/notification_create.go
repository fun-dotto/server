package repository

import (
	"context"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	notification.ID = uuid.New().String()

	dbNotification := database.NotificationFromDomain(notification)

	uniqueTargets := uniqueTargetUsers(notification.TargetUsers)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbNotification).Error; err != nil {
			return err
		}

		if len(uniqueTargets) > 0 {
			targets := make([]database.NotificationTargetUser, 0, len(uniqueTargets))
			for _, t := range uniqueTargets {
				targets = append(targets, database.NotificationTargetUser{
					NotificationID: notification.ID,
					UserID:         t.UserID,
					NotifiedAt:     t.NotifiedAt,
				})
			}
			if err := tx.Create(&targets).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Notification{}, err
	}

	return dbNotification.ToDomain(uniqueTargets), nil
}
