package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *NotificationRepository) GetNotificationsByIDs(ctx context.Context, ids []string) ([]domain.Notification, error) {
	uniqueIDs := uniqueStrings(ids)
	if len(uniqueIDs) == 0 {
		return []domain.Notification{}, nil
	}

	var dbNotifications []model.Notification
	if err := r.db.WithContext(ctx).Where("id IN ?", uniqueIDs).Find(&dbNotifications).Error; err != nil {
		return nil, err
	}
	if len(dbNotifications) == 0 {
		return []domain.Notification{}, nil
	}

	existingIDs := make([]string, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		existingIDs = append(existingIDs, n.ID)
	}

	var allTargets []model.NotificationTargetUser
	if err := r.db.WithContext(ctx).Where("notification_id IN ?", existingIDs).Find(&allTargets).Error; err != nil {
		return nil, err
	}

	targetMap := make(map[string][]domain.NotificationTargetUser)
	for _, t := range allTargets {
		key := t.NotificationID.String()
		targetMap[key] = append(targetMap[key], domain.NotificationTargetUser{
			UserID:     t.UserID,
			NotifiedAt: t.NotifiedAt,
		})
	}

	notifications := make([]domain.Notification, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		notifications = append(notifications, notificationToDomain(n, targetMap[n.ID]))
	}

	return notifications, nil
}

func (r *NotificationRepository) DispatchNotifications(ctx context.Context, deliveries map[string][]string) ([]domain.Notification, error) {
	if len(deliveries) == 0 {
		return []domain.Notification{}, nil
	}

	now := time.Now()
	notificationIDs := make([]string, 0, len(deliveries))
	for nid, userIDs := range deliveries {
		uniqueUsers := uniqueStrings(userIDs)
		if len(uniqueUsers) == 0 {
			continue
		}
		db := r.db.WithContext(ctx).Model(&model.NotificationTargetUser{}).
			Where("notification_id = ? AND user_id IN ?", nid, uniqueUsers).
			Update("notified_at", now)
		if db.Error != nil {
			return nil, db.Error
		}
		if db.RowsAffected > 0 {
			notificationIDs = append(notificationIDs, nid)
		}
	}

	if len(notificationIDs) == 0 {
		return []domain.Notification{}, nil
	}

	return r.GetNotificationsByIDs(ctx, notificationIDs)
}
