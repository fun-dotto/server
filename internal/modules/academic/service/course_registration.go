package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type courseRegistrationRepository interface {
	List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error)
	Create(ctx context.Context, cr domain.CourseRegistration) (domain.CourseRegistration, error)
	Delete(ctx context.Context, id string) error
}

type CourseRegistrationService struct {
	repo courseRegistrationRepository
}

func NewCourseRegistrationService(repo courseRegistrationRepository) *CourseRegistrationService {
	return &CourseRegistrationService{repo: repo}
}

func (s *CourseRegistrationService) List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error) {
	return s.repo.List(ctx, filter)
}

func (s *CourseRegistrationService) Create(ctx context.Context, cr domain.CourseRegistration) (domain.CourseRegistration, error) {
	return s.repo.Create(ctx, cr)
}

func (s *CourseRegistrationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
