package handler

import (
	"github.com/fun-dotto/user-api/internal/service"
)

type Handler struct {
	userService         *service.UserService
	fcmTokenService     *service.FCMTokenService
	notificationService *service.NotificationService
}

func NewHandler(userService *service.UserService, fcmTokenService *service.FCMTokenService, notificationService *service.NotificationService) *Handler {
	return &Handler{
		userService:         userService,
		fcmTokenService:     fcmTokenService,
		notificationService: notificationService,
	}
}
