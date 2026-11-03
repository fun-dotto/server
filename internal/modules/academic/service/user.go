package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type userRepository interface {
	FindByID(ctx context.Context, id string) (domain.User, bool, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

// FindByID はユーザーを取得する。見つからない場合は found=false を返す。
func (s *UserService) FindByID(ctx context.Context, id string) (domain.User, bool, error) {
	return s.repo.FindByID(ctx, id)
}
