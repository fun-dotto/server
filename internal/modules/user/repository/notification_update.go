package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

func (r *NotificationRepository) UpdateNotification(notification domain.Notification) (domain.Notification, error) {
	var dbNotification database.Notification

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing database.Notification
		if err := tx.First(&existing, "id = ?", notification.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		dbNotification = database.NotificationFromDomain(notification)

		if err := tx.Omit("IsNotified").Save(&dbNotification).Error; err != nil {
			return err
		}

		if err := tx.Where("notification_id = ?", notification.ID).Delete(&database.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		uniqueIDs := uniqueStrings(notification.TargetUserIDs)
		if len(uniqueIDs) > 0 {
			targets := make([]database.NotificationTargetUser, 0, len(uniqueIDs))
			for _, userID := range uniqueIDs {
				targets = append(targets, database.NotificationTargetUser{
					NotificationID: notification.ID,
					UserID:         userID,
				})
			}
			if err := tx.Create(&targets).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Notification{}, domain.ErrNotFound
		}
		return domain.Notification{}, err
	}

	return dbNotification.ToDomain(uniqueStrings(notification.TargetUserIDs)), nil
}
