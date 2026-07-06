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

	return r.hydrateNotifications(ctx, dbNotifications, false)
}
