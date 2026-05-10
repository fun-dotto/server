package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) TimetableItemsV1Create(ctx context.Context, request api.TimetableItemsV1CreateRequestObject) (api.TimetableItemsV1CreateResponseObject, error) {
	domainItem := toDomainTimetableItemFromRequest(*request.Body)
	created, err := h.timetableItemSvc.Create(ctx, domainItem)
	if err != nil {
		return nil, err
	}
	return api.TimetableItemsV1Create201JSONResponse{TimetableItem: timetableItemToAPI(created)}, nil
}
