package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

func (s *NotificationService) ListNotifications(ctx context.Context, filter domain.NotificationListFilter) ([]domain.Notification, error) {
	return s.repo.ListNotifications(ctx, filter)
}
