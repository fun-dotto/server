package repository

import "context"

// MarkUsersAsNotified は通知ID毎に指定ユーザーの notification_target_users.notified_at を現在時刻で更新する。
// 既に notified_at が入っているユーザーは上書きしない。
func (r *NotificationRepository) MarkUsersAsNotified(ctx context.Context, deliveries map[string][]string) error {
	_, err := r.markNotified(ctx, deliveries, true)
	return err
}
