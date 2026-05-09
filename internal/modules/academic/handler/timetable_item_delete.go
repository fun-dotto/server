package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) TimetableItemsV1Delete(ctx context.Context, request api.TimetableItemsV1DeleteRequestObject) (api.TimetableItemsV1DeleteResponseObject, error) {
	err := h.timetableItemSvc.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.TimetableItemsV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.TimetableItemsV1Delete204Response{}, nil
}
