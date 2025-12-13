package handler

import (
	"net/http"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	announcementService *service.AnnouncementService
}

func NewHandler(announcementService *service.AnnouncementService) *Handler {
	return &Handler{announcementService: announcementService}
}

func (h *Handler) AnnouncementsList(c *gin.Context, params api.AnnouncementsListParams) {
	ctx := c.Request.Context()
	announcementQuery := toDomainAnnouncementQuery(params)

	announcements, err := h.announcementService.GetAnnouncements(ctx, announcementQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	c.JSON(http.StatusOK, apiAnnouncements)
}
