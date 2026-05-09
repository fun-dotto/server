package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type facultyRepository interface {
	List(ctx context.Context, ids []string) ([]domain.Faculty, error)
	GetByID(ctx context.Context, id string) (domain.Faculty, error)
	Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type FacultyService struct {
	repo facultyRepository
}

func NewFacultyService(repo facultyRepository) *FacultyService {
	return &FacultyService{repo: repo}
}

func (s *FacultyService) List(ctx context.Context, ids []string) ([]domain.Faculty, error) {
	return s.repo.List(ctx, ids)
}

func (s *FacultyService) GetByID(ctx context.Context, id string) (domain.Faculty, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FacultyService) Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	return s.repo.Create(ctx, faculty)
}

func (s *FacultyService) Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	return s.repo.Update(ctx, faculty)
}

func (s *FacultyService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
