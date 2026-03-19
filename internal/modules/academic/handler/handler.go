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

type facultyService interface {
	List(ctx context.Context, ids []string) ([]domain.Faculty, error)
	GetByID(ctx context.Context, id string) (domain.Faculty, error)
	Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	subjectSvc subjectService
	facultySvc facultyService
}

func NewHandler(
	subjectSvc subjectService,
	facultySvc facultyService,
) *Handler {
	return &Handler{
		subjectSvc: subjectSvc,
		facultySvc: facultySvc,
	}
}
