package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultiesV1List(ctx context.Context, request api.FacultiesV1ListRequestObject) (api.FacultiesV1ListResponseObject, error) {
	var ids []string
	if request.Params.Ids != nil {
		ids = *request.Params.Ids
	}

	faculties, err := h.facultySvc.List(ctx, ids)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1List200JSONResponse{Faculties: facultiesToAPI(faculties)}, nil
}
