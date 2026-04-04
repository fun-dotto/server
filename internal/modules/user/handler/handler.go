package handler

import (
	"github.com/fun-dotto/api-template/internal/service"
)

type Handler struct {
	userService     *service.UserService
	fcmTokenService *service.FCMTokenService
}

func NewHandler(userService *service.UserService, fcmTokenService *service.FCMTokenService) *Handler {
	return &Handler{
		userService:     userService,
		fcmTokenService: fcmTokenService,
	}
}
