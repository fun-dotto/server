package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type timetableItemRepository interface {
	List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error)
	Create(ctx context.Context, item domain.TimetableItem) (domain.TimetableItem, error)
	Delete(ctx context.Context, id string) error
}

type TimetableItemService struct {
	repo timetableItemRepository
}

func NewTimetableItemService(repo timetableItemRepository) *TimetableItemService {
	return &TimetableItemService{repo: repo}
}

func (s *TimetableItemService) List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error) {
	return s.repo.List(ctx, filter)
}

func (s *TimetableItemService) Create(ctx context.Context, item domain.TimetableItem) (domain.TimetableItem, error) {
	return s.repo.Create(ctx, item)
}

func (s *TimetableItemService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
