package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) RoomChangesV1List(ctx context.Context, request api.RoomChangesV1ListRequestObject) (api.RoomChangesV1ListResponseObject, error) {
	filter := buildRoomChangeListFilter(request.Params)

	roomChanges, err := h.roomChangeSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	res, err := roomChangesToAPI(roomChanges)
	if err != nil {
		return nil, err
	}
	return api.RoomChangesV1List200JSONResponse{RoomChanges: res}, nil
}
