package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type subjectService interface {
	List(ctx context.Context) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	Create(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	Update(ctx context.Context, id string, subject domain.Subject) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
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
