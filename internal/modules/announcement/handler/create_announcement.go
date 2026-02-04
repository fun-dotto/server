package handler

import (
	"context"
	"time"

	"github.com/google/uuid"

	api "github.com/fun-dotto/announcement-api/generated"
)

func (h *Handler) AnnouncementsV1Create(ctx context.Context, request api.AnnouncementsV1CreateRequestObject) (api.AnnouncementsV1CreateResponseObject, error) {
	id := uuid.New().String()
	domainAnnouncement := toDomainAnnouncementFromRequest(id, *request.Body, time.Now())

	created, err := h.announcementService.CreateAnnouncement(ctx, domainAnnouncement)
	if err != nil {
		return nil, err
	}

	return api.AnnouncementsV1Create200JSONResponse{
		Announcement: toApiAnnouncement(created),
	}, nil
}
