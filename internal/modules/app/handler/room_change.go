package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) RoomChangesV1List(ctx context.Context, request api.RoomChangesV1ListRequestObject) (api.RoomChangesV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toRoomChangeQuery(request.Params)

	changes, err := h.academicService.GetRoomChanges(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get room changes: %w", err)
	}

	apiChanges := make([]api.RoomChange, len(changes))
	for i, c := range changes {
		apiChanges[i] = toApiRoomChange(c)
	}

	return api.RoomChangesV1List200JSONResponse{
		RoomChanges: apiChanges,
	}, nil
}

func toRoomChangeQuery(params api.RoomChangesV1ListParams) domain.RoomChangeQuery {
	query := domain.RoomChangeQuery{
		SubjectIDs: params.SubjectIds,
	}
	if params.From != nil {
		t := params.From.Time
		query.From = &t
	}
	if params.Until != nil {
		t := params.Until.Time
		query.Until = &t
	}
	return query
}

func toApiRoomChange(c domain.RoomChange) api.RoomChange {
	return api.RoomChange{
		Id:      c.ID,
		Date:    openapi_types.Date{Time: c.Date},
		Period:  api.DottoFoundationV1Period(c.Period),
		Subject: toApiSubjectSummary(c.Subject),
		OriginalRoom: api.Room{
			Id:    c.OriginalRoom.ID,
			Name:  c.OriginalRoom.Name,
			Floor: api.DottoFoundationV1Floor(c.OriginalRoom.Floor),
		},
		NewRoom: api.Room{
			Id:    c.NewRoom.ID,
			Name:  c.NewRoom.Name,
			Floor: api.DottoFoundationV1Floor(c.NewRoom.Floor),
		},
	}
}
