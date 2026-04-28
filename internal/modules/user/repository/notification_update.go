package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *NotificationRepository) UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	var dbNotification database.Notification
	uniqueTargets := uniqueTargetUsers(notification.TargetUsers)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing database.Notification
		if err := tx.First(&existing, "id = ?", notification.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		dbNotification = database.NotificationFromDomain(notification)

		if err := tx.Save(&dbNotification).Error; err != nil {
			return err
		}

		var existingTargets []database.NotificationTargetUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("notification_id = ?", notification.ID).
			Find(&existingTargets).Error; err != nil {
			return err
		}
		existingNotifiedAt := make(map[string]*time.Time, len(existingTargets))
		for _, t := range existingTargets {
			existingNotifiedAt[t.UserID] = t.NotifiedAt
		}

		if err := tx.Where("notification_id = ?", notification.ID).Delete(&database.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		if len(uniqueTargets) > 0 {
			targets := make([]database.NotificationTargetUser, 0, len(uniqueTargets))
			for i, t := range uniqueTargets {
				notifiedAt := t.NotifiedAt
				if notifiedAt == nil {
					if prev, ok := existingNotifiedAt[t.UserID]; ok {
						notifiedAt = prev
						uniqueTargets[i].NotifiedAt = prev
					}
				}
				targets = append(targets, database.NotificationTargetUser{
					NotificationID: notification.ID,
					UserID:         t.UserID,
					NotifiedAt:     notifiedAt,
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

	return dbNotification.ToDomain(uniqueTargets), nil
}
