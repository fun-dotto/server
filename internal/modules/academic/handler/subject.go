package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
)

func (h *Handler) SubjectsV1List(ctx context.Context, _ api.SubjectsV1ListRequestObject) (api.SubjectsV1ListResponseObject, error) {
	subjects, err := h.subjectSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1List200JSONResponse{Subjects: subjectsToAPI(subjects)}, nil
}

func (h *Handler) SubjectsV1Create(ctx context.Context, request api.SubjectsV1CreateRequestObject) (api.SubjectsV1CreateResponseObject, error) {
	subject, err := h.subjectSvc.Create(ctx, subjectRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1Create201JSONResponse{Subject: subjectToAPI(subject)}, nil
}

func (h *Handler) SubjectsV1Detail(ctx context.Context, request api.SubjectsV1DetailRequestObject) (api.SubjectsV1DetailResponseObject, error) {
	subject, err := h.subjectSvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1Detail200JSONResponse{Subject: subjectToAPI(subject)}, nil
}

func (h *Handler) SubjectsV1Update(ctx context.Context, request api.SubjectsV1UpdateRequestObject) (api.SubjectsV1UpdateResponseObject, error) {
	subject, err := h.subjectSvc.Update(ctx, request.Id, subjectRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1Update200JSONResponse{Subject: subjectToAPI(subject)}, nil
}

func (h *Handler) SubjectsV1Delete(ctx context.Context, request api.SubjectsV1DeleteRequestObject) (api.SubjectsV1DeleteResponseObject, error) {
	if err := h.subjectSvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.SubjectsV1Delete204Response{}, nil
}
