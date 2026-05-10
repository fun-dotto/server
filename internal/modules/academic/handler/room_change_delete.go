package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) RoomChangesV1Delete(ctx context.Context, request api.RoomChangesV1DeleteRequestObject) (api.RoomChangesV1DeleteResponseObject, error) {
	if err := h.roomChangeSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RoomChangesV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.RoomChangesV1Delete204Response{}, nil
}
