package handler

import (
	"time"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AnnouncementsList(c *gin.Context) {
	// TODO: 実装
	announcements := []api.Announcement{
		{
			Date:     time.Now(),
			Id:       "1",
			IsActive: true,
			Title:    "Announcement 1",
			Url:      "https://example.com",
		},
	}

	c.JSON(200, announcements)
}
