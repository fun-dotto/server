package service

import "github.com/fun-dotto/user-api/internal/domain"

func (s *NotificationService) CreateNotification(notification domain.Notification) (domain.Notification, error) {
	return s.repo.CreateNotification(notification)
}
