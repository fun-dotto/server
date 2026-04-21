package repository

import (
	"context"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (r *NotificationRepository) GetNotificationsByIDs(ctx context.Context, ids []string) ([]domain.Notification, error) {
	uniqueIDs := uniqueStrings(ids)
	if len(uniqueIDs) == 0 {
		return []domain.Notification{}, nil
	}

	var dbNotifications []database.Notification
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

	var allTargets []database.NotificationTargetUser
	if err := r.db.WithContext(ctx).Where("notification_id IN ?", existingIDs).Find(&allTargets).Error; err != nil {
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

func (r *NotificationRepository) DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error) {
	uniqueIDs := uniqueStrings(ids)
	if len(uniqueIDs) == 0 {
		return []domain.Notification{}, nil
	}

	if err := r.db.WithContext(ctx).Model(&database.Notification{}).
		Where("id IN ?", uniqueIDs).
		Update("is_notified", true).Error; err != nil {
		return nil, err
	}

	notifications, err := r.GetNotificationsByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, err
	}

	for i := range notifications {
		notifications[i].IsNotified = true
	}

	return notifications, nil
}
