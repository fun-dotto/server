package service

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
)

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.ListUsers(ctx)
}
