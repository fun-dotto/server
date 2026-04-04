package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultiesV1List(ctx context.Context, request api.FacultiesV1ListRequestObject) (api.FacultiesV1ListResponseObject, error) {
	// Serviceを変えてnilを渡さなくしてもいいかも
	faculties, err := h.facultySvc.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1List200JSONResponse{Faculties: facultiesToAPI(faculties)}, nil
}
