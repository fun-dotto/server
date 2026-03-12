package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
)

type subjectRepository interface {
	List(ctx context.Context) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	Create(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	Update(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
}

type SubjectService struct {
	repo subjectRepository
}

func NewSubjectService(repo subjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) List(ctx context.Context) ([]domain.Subject, error) {
	return s.repo.List(ctx)
}

func (s *SubjectService) GetByID(ctx context.Context, id string) (domain.Subject, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubjectService) Create(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	subject.ID = uuid.New().String()
	return s.repo.Create(ctx, subject)
}

func (s *SubjectService) Update(ctx context.Context, id string, subject domain.Subject) (domain.Subject, error) {
	subject.ID = id
	return s.repo.Update(ctx, subject)
}

func (s *SubjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
