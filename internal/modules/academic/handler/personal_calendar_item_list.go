package handler

import (
	"context"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func datesToTime(dates []openapi_types.Date) []time.Time {
	result := make([]time.Time, len(dates))
	for i, d := range dates {
		result[i] = d.Time
	}
	return result
}

func (h *Handler) PersonalCalendarItemsV1List(
	ctx context.Context,
	request api.PersonalCalendarItemsV1ListRequestObject) (api.PersonalCalendarItemsV1ListResponseObject, error) {
	dates := datesToTime(request.Params.Dates)
	items, err := h.personalCalendarItemSvc.List(ctx, request.Params.UserId, dates)
	if err != nil {
		return nil, err
	}
	return api.PersonalCalendarItemsV1List200JSONResponse{
		PersonalCalendarItems: personalCalendarItemsToAPI(items),
	}, nil
}
