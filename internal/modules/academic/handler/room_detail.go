package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) RoomsV1Detail(ctx context.Context, request api.RoomsV1DetailRequestObject) (api.RoomsV1DetailResponseObject, error) {
	room, err := h.roomSvc.GetByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RoomsV1Detail404Response{}, nil
		}
		return nil, err
	}
	return api.RoomsV1Detail200JSONResponse{Room: roomToAPI(room)}, nil
}
