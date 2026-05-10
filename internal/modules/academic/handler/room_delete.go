package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/repository"
	"gorm.io/gorm"
)

func (h *Handler) RoomsV1Delete(ctx context.Context, request api.RoomsV1DeleteRequestObject) (api.RoomsV1DeleteResponseObject, error) {
	if err := h.roomSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RoomsV1Delete404Response{}, nil
		}
		if errors.Is(err, repository.ErrRoomInUse) {
			return api.RoomsV1Delete409Response{}, nil
		}
		return nil, err
	}
	return api.RoomsV1Delete204Response{}, nil
}
