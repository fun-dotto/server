package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) SyllabusV1Detail(ctx context.Context, request api.SyllabusV1DetailRequestObject) (api.SyllabusV1DetailResponseObject, error) {
	syllabus, err := h.subjectSvc.GetSyllabus(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.SyllabusV1Detail404Response{}, nil
		}
		return nil, err
	}
	return api.SyllabusV1Detail200JSONResponse{Syllabus: syllabusToAPI(syllabus)}, nil
}
