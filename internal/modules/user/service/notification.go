package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

type NotificationRepository interface {
	ListNotifications(ctx context.Context, filter domain.NotificationListFilter) ([]domain.Notification, error)
	CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	DeleteNotification(ctx context.Context, id string) error
}

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}
