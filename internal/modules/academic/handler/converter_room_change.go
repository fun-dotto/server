package handler

import (
	"fmt"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func roomChangeToAPI(d domain.RoomChange) (api.RoomChange, error) {
	date, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		return api.RoomChange{}, fmt.Errorf("failed to parse room change date %q: %w", d.Date, err)
	}
	return api.RoomChange{
		Id:           d.ID,
		Subject:      subjectToAPI(d.Subject),
		Date:         openapi_types.Date{Time: date},
		Period:       api.DottoFoundationV1Period(d.Period),
		OriginalRoom: roomToAPI(d.OriginalRoom),
		NewRoom:      roomToAPI(d.NewRoom),
	}, nil
}

func roomChangesToAPI(ds []domain.RoomChange) ([]api.RoomChange, error) {
	result := make([]api.RoomChange, len(ds))
	for i, d := range ds {
		r, err := roomChangeToAPI(d)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func buildRoomChangeListFilter(params api.RoomChangesV1ListParams) domain.RoomChangeListFilter {
	filter := domain.RoomChangeListFilter{}
	if params.SubjectIds != nil {
		filter.SubjectIDs = *params.SubjectIds
	}
	if params.From != nil {
		filter.From = &params.From.Time
	}
	if params.Until != nil {
		filter.Until = &params.Until.Time
	}
	return filter
}

func toDomainRoomChangeFromRequest(req api.RoomChangeRequest) domain.RoomChange {
	return domain.RoomChange{
		Subject:      domain.Subject{ID: req.SubjectId},
		Date:         req.Date.Format("2006-01-02"),
		Period:       domain.Period(req.Period),
		OriginalRoom: domain.Room{ID: req.OriginalRoomId},
		NewRoom:      domain.Room{ID: req.NewRoomId},
	}
}
