package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

func (s *NotificationService) UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return s.repo.UpdateNotification(ctx, notification)
}
