package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) RoomsV1List(ctx context.Context, request api.RoomsV1ListRequestObject) (api.RoomsV1ListResponseObject, error) {
	filter := buildRoomListFilter(request.Params)

	rooms, err := h.roomSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return api.RoomsV1List200JSONResponse{Rooms: roomsToAPI(rooms)}, nil
}
