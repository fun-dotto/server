package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
)

func (s *NotificationService) UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return s.repo.UpdateNotification(ctx, notification)
}
