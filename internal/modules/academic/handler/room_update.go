package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) RoomsV1Update(ctx context.Context, request api.RoomsV1UpdateRequestObject) (api.RoomsV1UpdateResponseObject, error) {
	domainRoom := toDomainRoomFromRequest(request.Id, *request.Body)
	updated, err := h.roomSvc.Update(ctx, domainRoom)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RoomsV1Update404Response{}, nil
		}
		return nil, err
	}
	return api.RoomsV1Update200JSONResponse{Room: roomToAPI(updated)}, nil
}
