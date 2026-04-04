package handler

import (
	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func personalCalendarItemToAPI(d domain.PersonalCalendarItem) api.PersonalCalendarItem {
	rooms := make([]api.Room, len(d.Rooms))
	for i, r := range d.Rooms {
		rooms[i] = roomToAPI(r)
	}

	return api.PersonalCalendarItem{
		Date:    d.Date,
		Period:  api.DottoFoundationV1Period(d.Period),
		Subject: subjectToAPI(d.Subject),
		Rooms:   rooms,
		Status:  api.DottoFoundationV1PersonalCalendarItemStatus(d.Status),
	}
}

func personalCalendarItemsToAPI(ds []domain.PersonalCalendarItem) []api.PersonalCalendarItem {
	result := make([]api.PersonalCalendarItem, len(ds))
	for i, d := range ds {
		result[i] = personalCalendarItemToAPI(d)
	}
	return result
}
