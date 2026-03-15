package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type subjectService interface {
	List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	Upsert(ctx context.Context, syllabusID string) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
	GetSyllabus(ctx context.Context, subjectID string) (domain.Syllabus, error)
}

type Handler struct {
	subjectSvc subjectService
}

func NewHandler(
	subjectSvc subjectService,
) *Handler {
	return &Handler{
		subjectSvc: subjectSvc,
	}
}
