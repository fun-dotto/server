package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type roomRepository interface {
	List(ctx context.Context, filter domain.RoomListFilter) ([]domain.Room, error)
	GetByID(ctx context.Context, id string) (domain.Room, error)
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	Update(ctx context.Context, room domain.Room) (domain.Room, error)
	Delete(ctx context.Context, id string) error
}

type RoomService struct {
	repo roomRepository
}

func NewRoomService(repo roomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) List(ctx context.Context, filter domain.RoomListFilter) ([]domain.Room, error) {
	return s.repo.List(ctx, filter)
}

func (s *RoomService) GetByID(ctx context.Context, id string) (domain.Room, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoomService) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	return s.repo.Create(ctx, room)
}

func (s *RoomService) Update(ctx context.Context, room domain.Room) (domain.Room, error) {
	return s.repo.Update(ctx, room)
}

func (s *RoomService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
