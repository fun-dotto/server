package service

import "github.com/fun-dotto/user-api/internal/domain"

func (s *NotificationService) ListNotifications(filter domain.NotificationListFilter) ([]domain.Notification, error) {
	return s.repo.ListNotifications(filter)
}
