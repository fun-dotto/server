package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) MakeupClassesV1Delete(ctx context.Context, request api.MakeupClassesV1DeleteRequestObject) (api.MakeupClassesV1DeleteResponseObject, error) {
	if err := h.makeupClassSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.MakeupClassesV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.MakeupClassesV1Delete204Response{}, nil
}
