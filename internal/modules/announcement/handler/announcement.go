package handler

import (
	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct{}

func (h *AnnouncementHandler) AnnouncementsList(c *gin.Context) {
	// TODO: 実装
	announcements := []api.Announcement{}
	c.JSON(200, announcements)
}
