package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
)

func (s *NotificationService) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return s.repo.CreateNotification(ctx, notification)
}
