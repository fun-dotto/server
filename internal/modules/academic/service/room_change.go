package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type roomChangeRepository interface {
	List(ctx context.Context, filter domain.RoomChangeListFilter) ([]domain.RoomChange, error)
	GetByID(ctx context.Context, id string) (domain.RoomChange, error)
	Create(ctx context.Context, rc domain.RoomChange) (domain.RoomChange, error)
	Delete(ctx context.Context, id string) error
}

type RoomChangeService struct {
	repo roomChangeRepository
}

func NewRoomChangeService(repo roomChangeRepository) *RoomChangeService {
	return &RoomChangeService{repo: repo}
}

func (s *RoomChangeService) List(ctx context.Context, filter domain.RoomChangeListFilter) ([]domain.RoomChange, error) {
	return s.repo.List(ctx, filter)
}

func (s *RoomChangeService) GetByID(ctx context.Context, id string) (domain.RoomChange, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoomChangeService) Create(ctx context.Context, rc domain.RoomChange) (domain.RoomChange, error) {
	return s.repo.Create(ctx, rc)
}

func (s *RoomChangeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
