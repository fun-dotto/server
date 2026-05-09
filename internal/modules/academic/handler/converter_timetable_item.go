package handler

import (
	"strings"
	"time"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func buildTimetableItemListFilter(params api.TimetableItemsV1ListParams) domain.TimetableItemListFilter {
	filter := domain.TimetableItemListFilter{}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		// TODO: このデフォルト値設定のロジックは service 層に移すべき。
		// また、日本の大学の年度は4月始まりのため、1〜3月は前年度を返す必要がある。
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	semesters := make([]domain.CourseSemester, len(params.Semesters))
	for i, s := range params.Semesters {
		semesters[i] = domain.CourseSemester(s)
	}
	filter.Semesters = semesters
	return filter
}

func timetableSlotToAPI(slot domain.TimetableSlot) api.DottoFoundationV1TimetableSlot {
	return api.DottoFoundationV1TimetableSlot{
		DayOfWeek: api.DottoFoundationV1DayOfWeek(slot.DayOfWeek),
		Period:    api.DottoFoundationV1Period(slot.Period),
	}
}

func timetableItemToAPI(d domain.TimetableItem) api.TimetableItem {
	var slot *api.DottoFoundationV1TimetableSlot
	if d.Slot != nil &&
		strings.TrimSpace(string(d.Slot.DayOfWeek)) != "" &&
		strings.TrimSpace(string(d.Slot.Period)) != "" {
		s := timetableSlotToAPI(*d.Slot)
		slot = &s
	}

	rooms := make([]api.Room, len(d.Rooms))
	for i, r := range d.Rooms {
		rooms[i] = roomToAPI(r)
	}

	return api.TimetableItem{
		Id:      d.ID,
		Subject: subjectToAPI(d.Subject),
		Slot:    slot,
		Rooms:   rooms,
	}
}

func timetableItemsToAPI(ds []domain.TimetableItem) []api.TimetableItem {
	result := make([]api.TimetableItem, len(ds))
	for i, d := range ds {
		result[i] = timetableItemToAPI(d)
	}
	return result
}

func toDomainTimetableItemFromRequest(req api.TimetableItemRequest) domain.TimetableItem {
	item := domain.TimetableItem{
		Subject: domain.Subject{ID: req.SubjectId},
	}
	if req.Slot != nil {
		dow := strings.TrimSpace(string(req.Slot.DayOfWeek))
		per := strings.TrimSpace(string(req.Slot.Period))
		if dow != "" && per != "" {
			item.Slot = &domain.TimetableSlot{
				DayOfWeek: domain.DayOfWeek(dow),
				Period:    domain.Period(per),
			}
		}
	}
	rooms := make([]domain.Room, len(req.RoomIds))
	for i, id := range req.RoomIds {
		rooms[i] = domain.Room{ID: id}
	}
	item.Rooms = rooms
	return item
}
