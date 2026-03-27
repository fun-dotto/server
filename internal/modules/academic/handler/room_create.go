package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/google/uuid"
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
