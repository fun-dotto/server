package service

import (
	"context"

	"github.com/fun-dotto/subject-api/internal/domain"
	"github.com/google/uuid"
)

type courseRepository interface {
	List(ctx context.Context) ([]domain.Course, error)
	GetByID(ctx context.Context, id string) (domain.Course, error)
	Create(ctx context.Context, course domain.Course) (domain.Course, error)
	Update(ctx context.Context, course domain.Course) (domain.Course, error)
	Delete(ctx context.Context, id string) error
}

type CourseService struct {
	repo courseRepository
}

func NewCourseService(repo courseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) List(ctx context.Context) ([]domain.Course, error) {
	return s.repo.List(ctx)
}

func (s *CourseService) GetByID(ctx context.Context, id string) (domain.Course, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CourseService) Create(ctx context.Context, course domain.Course) (domain.Course, error) {
	course.ID = uuid.New().String()
	return s.repo.Create(ctx, course)
}

func (s *CourseService) Update(ctx context.Context, id string, course domain.Course) (domain.Course, error) {
	course.ID = id
	return s.repo.Update(ctx, course)
}

func (s *CourseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
