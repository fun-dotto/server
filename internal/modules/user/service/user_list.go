package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.ListUsers(ctx)
}
