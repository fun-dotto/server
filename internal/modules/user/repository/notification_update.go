package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *NotificationRepository) UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	notificationID, err := uuid.Parse(notification.ID)
	if err != nil {
		return domain.Notification{}, err
	}

	var dbNotification model.Notification
	uniqueTargets := uniqueTargetUsers(notification.TargetUsers)

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Notification
		if err := tx.First(&existing, "id = ?", notification.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		dbNotification = notificationFromDomain(notification)

		if err := tx.Save(&dbNotification).Error; err != nil {
			return err
		}

		var existingTargets []model.NotificationTargetUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("notification_id = ?", notification.ID).
			Find(&existingTargets).Error; err != nil {
			return err
		}
		existingNotifiedAt := make(map[string]*time.Time, len(existingTargets))
		for _, t := range existingTargets {
			existingNotifiedAt[t.UserID] = t.NotifiedAt
		}

		if err := tx.Where("notification_id = ?", notification.ID).Delete(&model.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		if len(uniqueTargets) > 0 {
			targets := make([]model.NotificationTargetUser, 0, len(uniqueTargets))
			for i, t := range uniqueTargets {
				notifiedAt := t.NotifiedAt
				if notifiedAt == nil {
					if prev, ok := existingNotifiedAt[t.UserID]; ok {
						notifiedAt = prev
						uniqueTargets[i].NotifiedAt = prev
					}
				}
				targets = append(targets, model.NotificationTargetUser{
					NotificationID: notificationID,
					UserID:         t.UserID,
					NotifiedAt:     notifiedAt,
				})
			}
			// Notification/User の関連は親側で更新済みなので Omit して GORM の auto-upsert を抑止する。
			if err := tx.Omit("Notification", "User").Create(&targets).Error; err != nil {
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

	return notificationToDomain(dbNotification, uniqueTargets), nil
}
