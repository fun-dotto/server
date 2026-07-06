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

	return r.hydrateNotifications(ctx, dbNotifications, false)
}

// markNotified は通知ID毎に指定ユーザーの notification_target_users.notified_at を現在時刻で更新する。
// onlyPending が true の場合、既に notified_at が入っているユーザーは上書きしない。
// updatedNotificationIDs には実際に1件以上更新された通知IDが返る。
func (r *NotificationRepository) markNotified(ctx context.Context, deliveries map[string][]string, onlyPending bool) (updatedNotificationIDs []string, err error) {
	if len(deliveries) == 0 {
		return nil, nil
	}

	now := time.Now()
	for nid, userIDs := range deliveries {
		uniqueUsers := uniqueStrings(userIDs)
		if len(uniqueUsers) == 0 {
			continue
		}
		query := r.db.WithContext(ctx).Model(&model.NotificationTargetUser{}).
			Where("notification_id = ? AND user_id IN ?", nid, uniqueUsers)
		if onlyPending {
			query = query.Where("notified_at IS NULL")
		}
		db := query.Update("notified_at", now)
		if db.Error != nil {
			return nil, db.Error
		}
		if db.RowsAffected > 0 {
			updatedNotificationIDs = append(updatedNotificationIDs, nid)
		}
	}
	return updatedNotificationIDs, nil
}

func (r *NotificationRepository) DispatchNotifications(ctx context.Context, deliveries map[string][]string) ([]domain.Notification, error) {
	notificationIDs, err := r.markNotified(ctx, deliveries, false)
	if err != nil {
		return nil, err
	}
	if len(notificationIDs) == 0 {
		return []domain.Notification{}, nil
	}
	return r.GetNotificationsByIDs(ctx, notificationIDs)
}
