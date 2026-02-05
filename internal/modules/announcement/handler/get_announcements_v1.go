package handler

import (
	"context"

	api "github.com/fun-dotto/announcement-api/generated"
)

func (h *Handler) AnnouncementsV1List(ctx context.Context, request api.AnnouncementsV1ListRequestObject) (api.AnnouncementsV1ListResponseObject, error) {
	announcementQuery := toDomainAnnouncementQuery(request.Params)

	announcements, err := h.announcementService.GetAnnouncements(ctx, announcementQuery)
	if err != nil {
		return nil, err
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	return api.AnnouncementsV1List200JSONResponse{
		Announcements: apiAnnouncements,
	}, nil
}
