package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) SyllabusV1Detail(ctx context.Context, request api.SyllabusV1DetailRequestObject) (api.SyllabusV1DetailResponseObject, error) {
	syllabus, err := h.subjectSvc.GetSyllabus(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.SyllabusV1Detail200JSONResponse{Syllabus: syllabusToAPI(syllabus)}, nil
}
