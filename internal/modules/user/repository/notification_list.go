package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *NotificationRepository) ListNotifications(ctx context.Context, filter domain.NotificationListFilter) ([]domain.Notification, error) {
	query := r.db.WithContext(ctx).Model(&model.Notification{})

	if filter.NotifyAtFrom != nil {
		query = query.Where("notify_before >= ?", *filter.NotifyAtFrom)
	}
	if filter.NotifyAtTo != nil {
		query = query.Where("notify_after <= ?", *filter.NotifyAtTo)
	}
	if filter.IsNotified != nil {
		if *filter.IsNotified {
			query = query.Where(`NOT EXISTS (SELECT 1 FROM notification_target_users tu WHERE tu.notification_id = notifications.id AND tu.notified_at IS NULL)
				AND EXISTS (SELECT 1 FROM notification_target_users tu WHERE tu.notification_id = notifications.id)`)
		} else {
			query = query.Where(`(EXISTS (SELECT 1 FROM notification_target_users tu WHERE tu.notification_id = notifications.id AND tu.notified_at IS NULL)
				OR NOT EXISTS (SELECT 1 FROM notification_target_users tu WHERE tu.notification_id = notifications.id))`)
		}
	}

	var dbNotifications []model.Notification
	if err := query.Order("notify_after DESC").Find(&dbNotifications).Error; err != nil {
		return nil, err
	}

	if len(dbNotifications) == 0 {
		return []domain.Notification{}, nil
	}

	notificationIDs := make([]string, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		notificationIDs = append(notificationIDs, n.ID)
	}

	var allTargets []model.NotificationTargetUser
	if err := r.db.WithContext(ctx).Where("notification_id IN ?", notificationIDs).Find(&allTargets).Error; err != nil {
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
