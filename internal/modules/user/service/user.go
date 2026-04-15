package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

type UserRepository interface {
	ListUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	UpsertUser(ctx context.Context, user domain.User) (domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpsertUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.UpsertUser(ctx, user)
}
