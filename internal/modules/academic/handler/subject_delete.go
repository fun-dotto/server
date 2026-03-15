package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) SubjectsV1Delete(ctx context.Context, request api.SubjectsV1DeleteRequestObject) (api.SubjectsV1DeleteResponseObject, error) {
	if err := h.subjectSvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.SubjectsV1Delete204Response{}, nil
}
