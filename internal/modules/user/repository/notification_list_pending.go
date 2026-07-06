package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *NotificationRepository) ListPendingNotifications(ctx context.Context, now time.Time) ([]domain.Notification, error) {
	var dbNotifications []model.Notification
	if err := r.db.WithContext(ctx).
		Where("notify_after <= ?", now).
		Where("notify_before > ?", now).
		Where(`EXISTS (
			SELECT 1 FROM notification_target_users tu
			WHERE tu.notification_id = notifications.id
			AND tu.notified_at IS NULL
		)`).
		Order("notify_after ASC").
		Find(&dbNotifications).Error; err != nil {
		return nil, err
	}
	return r.hydrateNotifications(ctx, dbNotifications, true)
}
