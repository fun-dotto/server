package handler

import (
	"context"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/service"
)

type Handler struct {
	announcementService *service.AnnouncementService
}

func NewHandler(announcementService *service.AnnouncementService) *Handler {
	return &Handler{announcementService: announcementService}
}

func (h *Handler) AnnouncementsList(ctx context.Context, request api.AnnouncementsListRequestObject) (api.AnnouncementsListResponseObject, error) {
	announcementQuery := toDomainAnnouncementQuery(request.Params)

	announcements, err := h.announcementService.GetAnnouncements(ctx, announcementQuery)
	if err != nil {
		return nil, err
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	return api.AnnouncementsList200JSONResponse(apiAnnouncements), nil
}
