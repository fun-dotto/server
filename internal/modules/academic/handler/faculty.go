package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
)

func (h *Handler) FacultiesV1List(ctx context.Context, _ api.FacultiesV1ListRequestObject) (api.FacultiesV1ListResponseObject, error) {
	faculties, err := h.facultySvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1List200JSONResponse{Faculties: facultiesToAPI(faculties)}, nil
}

func (h *Handler) FacultiesV1Create(ctx context.Context, request api.FacultiesV1CreateRequestObject) (api.FacultiesV1CreateResponseObject, error) {
	faculty, err := h.facultySvc.Create(ctx, facultyRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1Create201JSONResponse{Faculty: facultyToAPI(faculty)}, nil
}

func (h *Handler) FacultiesV1Detail(ctx context.Context, request api.FacultiesV1DetailRequestObject) (api.FacultiesV1DetailResponseObject, error) {
	faculty, err := h.facultySvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1Detail200JSONResponse{Faculty: facultyToAPI(faculty)}, nil
}

func (h *Handler) FacultiesV1Update(ctx context.Context, request api.FacultiesV1UpdateRequestObject) (api.FacultiesV1UpdateResponseObject, error) {
	faculty, err := h.facultySvc.Update(ctx, request.Id, facultyRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.FacultiesV1Update200JSONResponse{Faculty: facultyToAPI(faculty)}, nil
}

func (h *Handler) FacultiesV1Delete(ctx context.Context, request api.FacultiesV1DeleteRequestObject) (api.FacultiesV1DeleteResponseObject, error) {
	if err := h.facultySvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.FacultiesV1Delete204Response{}, nil
}
