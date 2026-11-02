package external

import (
	"time"

	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainTimetableItem は外部APIのTimetableItemをDomainのTimetableItemに変換する
func ToDomainTimetableItem(m academic_api.TimetableItem) domain.TimetableItem {
	var slot *domain.TimetableSlot
	if m.Slot != nil {
		slot = &domain.TimetableSlot{
			DayOfWeek: domain.DayOfWeek(m.Slot.DayOfWeek),
			Period:    domain.Period(m.Slot.Period),
		}
	}

	return domain.TimetableItem{
		ID:      m.Id,
		Slot:    slot,
		Rooms:   toDomainRooms(m.Rooms),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalTimetableItemQuery はDomainのTimetableItemQueryを外部APIのTimetableItemsV1ListParamsに変換する
func ToExternalTimetableItemQuery(q domain.TimetableItemQuery) *academic_api.TimetableItemsV1ListParams {
	params := &academic_api.TimetableItemsV1ListParams{
		Semesters: toExternalSemesters(q.Semesters),
		Year:      q.Year,
	}

	return params
}

// ToDomainPersonalCalendarItem は外部APIのPersonalCalendarItemをDomainに変換する
func ToDomainPersonalCalendarItem(m academic_api.PersonalCalendarItem) domain.PersonalCalendarItem {
	return domain.PersonalCalendarItem{
		Date:    m.Date.Time,
		Period:  domain.Period(m.Period),
		Rooms:   toDomainRooms(m.Rooms),
		Status:  domain.PersonalCalendarItemStatus(m.Status),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalPersonalCalendarItemParams は domain の検索条件を外部APIのクエリに変換する
func ToExternalPersonalCalendarItemParams(userID string, dates []time.Time) *academic_api.PersonalCalendarItemsV1ListParams {
	apiDates := make([]openapi_types.Date, len(dates))
	for i, d := range dates {
		apiDates[i] = openapi_types.Date{Time: d}
	}
	return &academic_api.PersonalCalendarItemsV1ListParams{
		UserId: userID,
		Dates:  apiDates,
	}
}

func toDomainRooms(rooms []academic_api.Room) []domain.Room {
	result := make([]domain.Room, len(rooms))
	for i, room := range rooms {
		result[i] = domain.Room{
			ID:    room.Id,
			Name:  room.Name,
			Floor: domain.Floor(room.Floor),
		}
	}
	return result
}
