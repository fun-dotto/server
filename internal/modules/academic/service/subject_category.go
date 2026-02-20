package service

import (
	"context"

	"github.com/fun-dotto/subject-api/internal/domain"
	"github.com/google/uuid"
)

type subjectCategoryRepository interface {
	List(ctx context.Context) ([]domain.SubjectCategory, error)
	GetByID(ctx context.Context, id string) (domain.SubjectCategory, error)
	Create(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error)
	Update(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error)
	Delete(ctx context.Context, id string) error
}

type SubjectCategoryService struct {
	repo subjectCategoryRepository
}

func NewSubjectCategoryService(repo subjectCategoryRepository) *SubjectCategoryService {
	return &SubjectCategoryService{repo: repo}
}

func (s *SubjectCategoryService) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	return s.repo.List(ctx)
}

func (s *SubjectCategoryService) GetByID(ctx context.Context, id string) (domain.SubjectCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubjectCategoryService) Create(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error) {
	category.ID = uuid.New().String()
	return s.repo.Create(ctx, category)
}

func (s *SubjectCategoryService) Update(ctx context.Context, id string, category domain.SubjectCategory) (domain.SubjectCategory, error) {
	category.ID = id
	return s.repo.Update(ctx, category)
}

func (s *SubjectCategoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
