package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) SubjectsV1Upsert(ctx context.Context, request api.SubjectsV1UpsertRequestObject) (api.SubjectsV1UpsertResponseObject, error) {
	subject, err := h.subjectSvc.Upsert(ctx, request.Body.SyllabusId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.SubjectsV1Upsert400Response{}, nil
		}
		return nil, err
	}
	return api.SubjectsV1Upsert200JSONResponse{Subject: subjectToAPI(subject)}, nil
}
