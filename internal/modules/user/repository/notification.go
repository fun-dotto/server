package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// hydrateNotifications はターゲットユーザーを取得して dbNotifications に紐付ける。
// onlyPending が true の場合、notified_at IS NULL のターゲットのみを対象とする。
func (r *NotificationRepository) hydrateNotifications(ctx context.Context, dbNotifications []model.Notification, onlyPending bool) ([]domain.Notification, error) {
	if len(dbNotifications) == 0 {
		return []domain.Notification{}, nil
	}

	notificationIDs := make([]string, 0, len(dbNotifications))
	for _, n := range dbNotifications {
		notificationIDs = append(notificationIDs, n.ID)
	}

	query := r.db.WithContext(ctx).Where("notification_id IN ?", notificationIDs)
	if onlyPending {
		query = query.Where("notified_at IS NULL")
	}

	var allTargets []model.NotificationTargetUser
	if err := query.Find(&allTargets).Error; err != nil {
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
