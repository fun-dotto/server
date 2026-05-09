package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) RoomChangesV1Create(ctx context.Context, request api.RoomChangesV1CreateRequestObject) (api.RoomChangesV1CreateResponseObject, error) {
	domainRC := toDomainRoomChangeFromRequest(*request.Body)
	created, err := h.roomChangeSvc.Create(ctx, domainRC)
	if err != nil {
		return nil, err
	}
	res, err := roomChangeToAPI(created)
	if err != nil {
		return nil, err
	}
	return api.RoomChangesV1Create201JSONResponse{RoomChange: res}, nil
}
