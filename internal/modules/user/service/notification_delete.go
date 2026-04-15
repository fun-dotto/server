package service

import "context"

func (s *NotificationService) DeleteNotification(ctx context.Context, id string) error {
	return s.repo.DeleteNotification(ctx, id)
}
