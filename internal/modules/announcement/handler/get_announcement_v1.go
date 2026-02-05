package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
)

func (h *Handler) AnnouncementsV1Detail(ctx context.Context, request api.AnnouncementsV1DetailRequestObject) (api.AnnouncementsV1DetailResponseObject, error) {
	announcement, err := h.announcementService.GetAnnouncementByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err // TODO: 404レスポンスを返すべき
		}
		return nil, err
	}

	return api.AnnouncementsV1Detail200JSONResponse{
		Announcement: toApiAnnouncement(announcement),
	}, nil
}
