package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) TimetableItemsV1List(ctx context.Context, request api.TimetableItemsV1ListRequestObject) (api.TimetableItemsV1ListResponseObject, error) {
	filter := buildTimetableItemListFilter(request.Params)

	items, err := h.timetableItemSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return api.TimetableItemsV1List200JSONResponse{TimetableItems: timetableItemsToAPI(items)}, nil
}
