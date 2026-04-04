package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) PersonalCalendarItemsV1List(
	ctx context.Context,
	request api.PersonalCalendarItemsV1ListRequestObject) (api.PersonalCalendarItemsV1ListResponseObject, error) {
	items, err := h.personalCalendarItemSvc.List(ctx, request.Params.UserId, request.Params.Dates)
	if err != nil {
		return nil, err
	}
	return api.PersonalCalendarItemsV1List200JSONResponse{
		PersonalCalendarItems: personalCalendarItemsToAPI(items),
	}, nil
}
