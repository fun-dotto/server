package handler

import (
	"github.com/fun-dotto/server/internal/modules/announcement/service"
)

type Handler struct {
	announcementService *service.AnnouncementService
}

func NewHandler(announcementService *service.AnnouncementService) *Handler {
	return &Handler{announcementService: announcementService}
}
