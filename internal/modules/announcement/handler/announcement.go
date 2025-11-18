package handler

import (
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	announcementService *service.AnnouncementService
}

func NewHandler(announcementService *service.AnnouncementService) *Handler {
	return &Handler{announcementService: announcementService}
}

func (h *Handler) AnnouncementsList(c *gin.Context) {
	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, announcements)
}
