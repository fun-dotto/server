package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) SubjectsV1Detail(ctx context.Context, request api.SubjectsV1DetailRequestObject) (api.SubjectsV1DetailResponseObject, error) {
	subject, err := h.subjectSvc.GetByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.SubjectsV1Detail404Response{}, nil
		}
		return nil, err
	}
	return api.SubjectsV1Detail200JSONResponse{Subject: subjectToDetailAPI(subject)}, nil
}
