package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) FacultiesV1Update(ctx context.Context, request api.FacultiesV1UpdateRequestObject) (api.FacultiesV1UpdateResponseObject, error) {
	domainFaculty := toDomainFacultyFromRequest(request.Id, *request.Body)
	updated, err := h.facultySvc.Update(ctx, domainFaculty)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.FacultiesV1Update404Response{}, nil
		}
		return nil, err
	}
	return api.FacultiesV1Update200JSONResponse{Faculty: facultyToAPI(updated)}, nil
}
