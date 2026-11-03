package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) TimetableItemsV1List(ctx context.Context, request api.TimetableItemsV1ListRequestObject) (api.TimetableItemsV1ListResponseObject, error) {
	filter := buildTimetableItemListFilter(request.Params)

	sort, err := h.buildSubjectListSort(ctx, request.Params.UserId, request.Params.XFlags)
	if err != nil {
		return nil, err
	}

	items, err := h.timetableItemSvc.List(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	return api.TimetableItemsV1List200JSONResponse{TimetableItems: timetableItemsToAPI(items)}, nil
}
