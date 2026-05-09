package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) FacultiesV1Delete(ctx context.Context, request api.FacultiesV1DeleteRequestObject) (api.FacultiesV1DeleteResponseObject, error) {
	if err := h.facultySvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.FacultiesV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.FacultiesV1Delete204Response{}, nil
}
