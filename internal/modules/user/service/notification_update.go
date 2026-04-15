package service

import "github.com/fun-dotto/user-api/internal/domain"

func (s *NotificationService) UpdateNotification(notification domain.Notification) (domain.Notification, error) {
	return s.repo.UpdateNotification(notification)
}
