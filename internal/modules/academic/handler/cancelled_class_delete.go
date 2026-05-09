package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) CancelledClassesV1Delete(ctx context.Context, request api.CancelledClassesV1DeleteRequestObject) (api.CancelledClassesV1DeleteResponseObject, error) {
	if err := h.cancelledClassSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.CancelledClassesV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.CancelledClassesV1Delete204Response{}, nil
}
