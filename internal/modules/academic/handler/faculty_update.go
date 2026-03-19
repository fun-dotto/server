package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultiesV1Update(ctx context.Context, request api.FacultiesV1UpdateRequestObject) (api.FacultiesV1UpdateResponseObject, error) {
	domainFaculty := toDomainFacultyFromRequest(request.Id, *request.Body)
	updated, err := h.facultySvc.Update(ctx, domainFaculty)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1Update200JSONResponse{Faculty: facultyToAPI(updated)}, nil
}
