package service

import "github.com/fun-dotto/user-api/internal/domain"

type NotificationRepository interface {
	ListNotifications(filter domain.NotificationListFilter) ([]domain.Notification, error)
	CreateNotification(notification domain.Notification) (domain.Notification, error)
	UpdateNotification(notification domain.Notification) (domain.Notification, error)
	DeleteNotification(id string) error
}

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}
