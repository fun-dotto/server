package handler

import (
	"context"
	"uuid"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) RoomsV1Create(ctx context.Context, request api.RoomsV1CreateRequestObject) (api.RoomsV1CreateResponseObject, error) {
	id := uuid.New().String()
	domainRoom := toDomainRoomFromRequest(id, *request.Body)
	created, err := h.roomSvc.Create(ctx, domainRoom)
	if err != nil {
		return nil, err
	}
	return api.RoomsV1Create201JSONResponse{Room: roomToAPI(created)}, nil
}
