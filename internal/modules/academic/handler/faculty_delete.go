package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultiesV1Delete(ctx context.Context, request api.FacultiesV1DeleteRequestObject) (api.FacultiesV1DeleteResponseObject, error) {
	if err := h.facultySvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.FacultiesV1Delete204Response{}, nil
}
