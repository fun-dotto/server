package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) SubjectsV1Delete(ctx context.Context, request api.SubjectsV1DeleteRequestObject) (api.SubjectsV1DeleteResponseObject, error) {
	if err := h.subjectSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.SubjectsV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.SubjectsV1Delete204Response{}, nil
}
