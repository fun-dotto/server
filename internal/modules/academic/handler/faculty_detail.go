package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"gorm.io/gorm"
)

func (h *Handler) FacultiesV1Detail(ctx context.Context, request api.FacultiesV1DetailRequestObject) (api.FacultiesV1DetailResponseObject, error) {
	faculty, err := h.facultySvc.GetByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.FacultiesV1Detail404Response{}, nil
		}
		return nil, err
	}
	return api.FacultiesV1Detail200JSONResponse{Faculty: facultyToAPI(faculty)}, nil
}
