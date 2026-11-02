package handler

import (
	"context"
	"uuid"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) FacultiesV1Create(ctx context.Context, request api.FacultiesV1CreateRequestObject) (api.FacultiesV1CreateResponseObject, error) {
	id := uuid.New().String()
	domainFaculty := toDomainFacultyFromRequest(id, *request.Body)
	created, err := h.facultySvc.Create(ctx, domainFaculty)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1Create201JSONResponse{Faculty: facultyToAPI(created)}, nil
}
