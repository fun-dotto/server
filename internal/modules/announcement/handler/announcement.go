package handler

import (
	"net/http"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
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
	announcements, err := h.announcementService.GetAnnouncements(params.IsActive)
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

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:       announcement.ID,
		Title:    announcement.Title,
		Date:     announcement.Date,
		Url:      announcement.URL,
		IsActive: announcement.IsActive,
	}
}
