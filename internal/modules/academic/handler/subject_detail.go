package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) SubjectsV1Detail(ctx context.Context, request api.SubjectsV1DetailRequestObject) (api.SubjectsV1DetailResponseObject, error) {
	subject, err := h.subjectSvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1Detail200JSONResponse{Subject: subjectToAPI(subject)}, nil
}
