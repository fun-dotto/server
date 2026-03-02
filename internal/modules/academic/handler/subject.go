package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/subject-api/generated"
)

// TODO: クエリパラメータ (q, grade, courses, class, classification, semester, requirementType, culturalSubjectCategory) を使ったフィルタリングを実装する
func (h *Handler) SubjectsV1List(ctx context.Context, _ api.SubjectsV1ListRequestObject) (api.SubjectsV1ListResponseObject, error) {
	subjects, err := h.subjectSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1List200JSONResponse{Subjects: subjectsToSummaryAPI(subjects)}, nil
}

// TODO: 実装する
func (h *Handler) SubjectsV1Upsert(_ context.Context, _ api.SubjectsV1UpsertRequestObject) (api.SubjectsV1UpsertResponseObject, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) SubjectsV1Detail(ctx context.Context, request api.SubjectsV1DetailRequestObject) (api.SubjectsV1DetailResponseObject, error) {
	subject, err := h.subjectSvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1Detail200JSONResponse{Subject: subjectToAPI(subject)}, nil
}

func (h *Handler) SubjectsV1Delete(ctx context.Context, request api.SubjectsV1DeleteRequestObject) (api.SubjectsV1DeleteResponseObject, error) {
	if err := h.subjectSvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.SubjectsV1Delete204Response{}, nil
}

// TODO: 実装する
func (h *Handler) SyllabusV1Detail(_ context.Context, _ api.SyllabusV1DetailRequestObject) (api.SyllabusV1DetailResponseObject, error) {
	return nil, errors.New("not implemented")
}
