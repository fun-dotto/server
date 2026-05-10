package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *NotificationRepository) UpsertNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	if notification.ID == "" {
		return domain.Notification{}, errors.New("notification ID is required for upsert")
	}

	dbNotification, err := notificationFromDomain(notification)
	if err != nil {
		return domain.Notification{}, err
	}
	notificationID, err := uuid.Parse(dbNotification.ID)
	if err != nil {
		return domain.Notification{}, err
	}
	uniqueTargets := uniqueTargetUsers(notification.TargetUsers)

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ID 衝突時は本文を更新しない (再通知・重複配信を防ぐため)。target_users の増減のみ下で同期する。
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&dbNotification).Error; err != nil {
			return err
		}

		if len(uniqueTargets) == 0 {
			return tx.Where("notification_id = ?", dbNotification.ID).
				Delete(&model.NotificationTargetUser{}).Error
		}

		userIDs := make([]string, 0, len(uniqueTargets))
		for _, t := range uniqueTargets {
			userIDs = append(userIDs, t.UserID)
		}
		if err := tx.Where("notification_id = ? AND user_id NOT IN ?", dbNotification.ID, userIDs).
			Delete(&model.NotificationTargetUser{}).Error; err != nil {
			return err
		}

		targets := make([]model.NotificationTargetUser, 0, len(uniqueTargets))
		for _, t := range uniqueTargets {
			targets = append(targets, model.NotificationTargetUser{
				NotificationID: notificationID,
				UserID:         t.UserID,
				NotifiedAt:     t.NotifiedAt,
			})
		}
		// 既存行は notified_at を保持したいので競合時は何もしない。
		// Notification/User の関連は親側で作成済みなので Omit して GORM の auto-upsert を抑止する。
		return tx.Omit("Notification", "User").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "notification_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).Create(&targets).Error
	})
	if err != nil {
		return domain.Notification{}, err
	}

	return notificationToDomain(&dbNotification, uniqueTargets), nil
}
