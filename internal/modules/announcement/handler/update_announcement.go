package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
)

func (h *Handler) AnnouncementsV1Update(ctx context.Context, request api.AnnouncementsV1UpdateRequestObject) (api.AnnouncementsV1UpdateResponseObject, error) {
	domainAnnouncement := toDomainAnnouncementFromRequest(request.Id, *request.Body)

	updated, err := h.announcementService.UpdateAnnouncement(ctx, domainAnnouncement)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err // TODO: 404レスポンスを返すべき
		}
		return nil, err
	}

	return api.AnnouncementsV1Update200JSONResponse{
		Announcement: toApiAnnouncement(updated),
	}, nil
}
