package service

import (
	"context"

	"github.com/fun-dotto/subject-api/internal/domain"
	"github.com/google/uuid"
)

type facultyRepository interface {
	List(ctx context.Context) ([]domain.Faculty, error)
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

func (s *FacultyService) List(ctx context.Context) ([]domain.Faculty, error) {
	return s.repo.List(ctx)
}

func (s *FacultyService) GetByID(ctx context.Context, id string) (domain.Faculty, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FacultyService) Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	faculty.ID = uuid.New().String()
	return s.repo.Create(ctx, faculty)
}

func (s *FacultyService) Update(ctx context.Context, id string, faculty domain.Faculty) (domain.Faculty, error) {
	faculty.ID = id
	return s.repo.Update(ctx, faculty)
}

func (s *FacultyService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
