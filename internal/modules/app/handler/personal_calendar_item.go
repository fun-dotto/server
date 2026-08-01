package handler

import (
	"context"
	"fmt"
	"time"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/middleware"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// PersonalCalendarItemsV1List 個人カレンダーアイテム一覧を取得する
func (h *Handler) PersonalCalendarItemsV1List(ctx context.Context, request api.PersonalCalendarItemsV1ListRequestObject) (api.PersonalCalendarItemsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return api.PersonalCalendarItemsV1List401Response{}, nil
	}

	dates := make([]time.Time, len(request.Params.Dates))
	for i, d := range request.Params.Dates {
		dates[i] = d.Time
	}

	items, err := h.academicService.GetPersonalCalendarItems(userID, dates)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal calendar items: %w", err)
	}

	apiItems := make([]api.PersonalCalendarItem, len(items))
	for i, item := range items {
		apiItems[i] = toAPIPersonalCalendarItem(item)
	}

	return api.PersonalCalendarItemsV1List200JSONResponse{
		PersonalCalendarItems: apiItems,
	}, nil
}

func toAPIPersonalCalendarItem(item domain.PersonalCalendarItem) api.PersonalCalendarItem {
	return api.PersonalCalendarItem{
		Date:    openapi_types.Date{Time: item.Date},
		Period:  api.DottoFoundationV1Period(item.Period),
		Rooms:   toApiRooms(item.Rooms),
		Status:  api.DottoFoundationV1PersonalCalendarItemStatus(item.Status),
		Subject: toApiSubjectSummary(item.Subject),
	}
}
