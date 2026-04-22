package repository

import (
	"context"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (r *NotificationRepository) ListNotifications(ctx context.Context, filter domain.NotificationListFilter) ([]domain.Notification, error) {
	query := r.db.WithContext(ctx).Model(&database.Notification{})

	if filter.NotifyAtFrom != nil {
		query = query.Where("notify_before >= ?", *filter.NotifyAtFrom)
	}
	if filter.NotifyAtTo != nil {
		query = query.Where("notify_after <= ?", *filter.NotifyAtTo)
	}
	if filter.IsNotified != nil {
		query = query.Where("is_notified = ?", *filter.IsNotified)
	}

	var dbNotifications []database.Notification
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

	var allTargets []database.NotificationTargetUser
	if err := r.db.WithContext(ctx).Where("notification_id IN ?", notificationIDs).Find(&allTargets).Error; err != nil {
		return nil, err
	}

	targetMap := make(map[string][]string)
	for _, t := range allTargets {
		targetMap[t.NotificationID] = append(targetMap[t.NotificationID], t.UserID)
	}

	notifications := make([]domain.Notification, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		notifications = append(notifications, n.ToDomain(targetMap[n.ID]))
	}

	return notifications, nil
}
