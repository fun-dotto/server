package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type makeupClassRepository interface {
	List(ctx context.Context, filter domain.MakeupClassListFilter) ([]domain.MakeupClass, error)
	GetByID(ctx context.Context, id string) (domain.MakeupClass, error)
	Create(ctx context.Context, mc domain.MakeupClass) (domain.MakeupClass, error)
	Delete(ctx context.Context, id string) error
}

type MakeupClassService struct {
	repo makeupClassRepository
}

func NewMakeupClassService(repo makeupClassRepository) *MakeupClassService {
	return &MakeupClassService{repo: repo}
}

func (s *MakeupClassService) List(ctx context.Context, filter domain.MakeupClassListFilter) ([]domain.MakeupClass, error) {
	return s.repo.List(ctx, filter)
}

func (s *MakeupClassService) GetByID(ctx context.Context, id string) (domain.MakeupClass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MakeupClassService) Create(ctx context.Context, mc domain.MakeupClass) (domain.MakeupClass, error) {
	return s.repo.Create(ctx, mc)
}

func (s *MakeupClassService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
