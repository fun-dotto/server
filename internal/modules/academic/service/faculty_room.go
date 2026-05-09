package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type facultyRoomRepository interface {
	List(ctx context.Context, filter domain.FacultyRoomListFilter) ([]domain.FacultyRoom, error)
	Create(ctx context.Context, fr domain.FacultyRoom) (domain.FacultyRoom, error)
	Delete(ctx context.Context, id string) error
}

type FacultyRoomService struct {
	repo facultyRoomRepository
}

func NewFacultyRoomService(repo facultyRoomRepository) *FacultyRoomService {
	return &FacultyRoomService{repo: repo}
}

func (s *FacultyRoomService) List(ctx context.Context, filter domain.FacultyRoomListFilter) ([]domain.FacultyRoom, error) {
	return s.repo.List(ctx, filter)
}

func (s *FacultyRoomService) Create(ctx context.Context, fr domain.FacultyRoom) (domain.FacultyRoom, error) {
	return s.repo.Create(ctx, fr)
}

func (s *FacultyRoomService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
