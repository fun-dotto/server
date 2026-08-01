package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// TimetableItemsV1List 時間割アイテム一覧を取得する
func (h *Handler) TimetableItemsV1List(ctx context.Context, request api.TimetableItemsV1ListRequestObject) (api.TimetableItemsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toTimetableItemQuery(request.Params)

	items, err := h.academicService.GetTimetableItems(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get timetable items: %w", err)
	}

	apiItems := make([]api.TimetableItem, len(items))
	for i, item := range items {
		apiItems[i] = toApiTimetableItem(item)
	}

	return api.TimetableItemsV1List200JSONResponse{
		TimetableItems: apiItems,
	}, nil
}

// toTimetableItemQuery は TimetableItemsV1List の API パラメータを domain.TimetableItemQuery に変換する
func toTimetableItemQuery(params api.TimetableItemsV1ListParams) domain.TimetableItemQuery {
	query := domain.TimetableItemQuery{Year: params.Year}
	query.Semesters = make([]domain.CourseSemester, len(params.Semesters))
	for i, semester := range params.Semesters {
		query.Semesters[i] = domain.CourseSemester(semester)
	}
	return query
}

// toApiTimetableItem はDomainの時間割アイテムをAPIの時間割アイテムに変換する
func toApiTimetableItem(item domain.TimetableItem) api.TimetableItem {
	var slot *api.DottoFoundationV1TimetableSlot
	if item.Slot != nil {
		slot = &api.DottoFoundationV1TimetableSlot{
			DayOfWeek: api.DottoFoundationV1DayOfWeek(item.Slot.DayOfWeek),
			Period:    api.DottoFoundationV1Period(item.Slot.Period),
		}
	}

	return api.TimetableItem{
		Id:      item.ID,
		Slot:    slot,
		Rooms:   toApiRooms(item.Rooms),
		Subject: toApiSubjectSummary(item.Subject),
	}
}

func toApiRooms(rooms []domain.Room) []api.Room {
	result := make([]api.Room, len(rooms))
	for i, room := range rooms {
		result[i] = api.Room{
			Id:    room.ID,
			Name:  room.Name,
			Floor: api.DottoFoundationV1Floor(room.Floor),
		}
	}
	return result
}
