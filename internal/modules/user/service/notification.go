package service

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/fun-dotto/user-api/internal/domain"
)

type NotificationRepository interface {
	ListNotifications(ctx context.Context, filter domain.NotificationListFilter) ([]domain.Notification, error)
	CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	DeleteNotification(ctx context.Context, id string) error
	GetNotificationsByIDs(ctx context.Context, ids []string) ([]domain.Notification, error)
	DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error)
}

type FCMTokenRepositoryForNotification interface {
	ListFCMTokens(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error)
}

type NotificationService struct {
	repo             NotificationRepository
	fcmTokenRepo     FCMTokenRepositoryForNotification
	messagingClient  *messaging.Client
}

func NewNotificationService(
	repo NotificationRepository,
	fcmTokenRepo FCMTokenRepositoryForNotification,
	messagingClient *messaging.Client,
) *NotificationService {
	return &NotificationService{
		repo:            repo,
		fcmTokenRepo:    fcmTokenRepo,
		messagingClient: messagingClient,
	}
}
