package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) SubjectsV1List(ctx context.Context, request api.SubjectsV1ListRequestObject) (api.SubjectsV1ListResponseObject, error) {
	filter := buildSubjectListFilter(request.Params)

	subjects, err := h.subjectSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1List200JSONResponse{Subjects: subjectsToAPI(subjects)}, nil
}
