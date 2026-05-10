package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type cancelledClassRepository interface {
	List(ctx context.Context, filter domain.CancelledClassListFilter) ([]domain.CancelledClass, error)
	GetByID(ctx context.Context, id string) (domain.CancelledClass, error)
	Create(ctx context.Context, cc domain.CancelledClass) (domain.CancelledClass, error)
	Delete(ctx context.Context, id string) error
}

type CancelledClassService struct {
	repo cancelledClassRepository
}

func NewCancelledClassService(repo cancelledClassRepository) *CancelledClassService {
	return &CancelledClassService{repo: repo}
}

func (s *CancelledClassService) List(ctx context.Context, filter domain.CancelledClassListFilter) ([]domain.CancelledClass, error) {
	return s.repo.List(ctx, filter)
}

func (s *CancelledClassService) GetByID(ctx context.Context, id string) (domain.CancelledClass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CancelledClassService) Create(ctx context.Context, cc domain.CancelledClass) (domain.CancelledClass, error) {
	return s.repo.Create(ctx, cc)
}

func (s *CancelledClassService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
