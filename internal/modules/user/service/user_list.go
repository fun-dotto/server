package service

import "github.com/fun-dotto/user-api/internal/domain"

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}
