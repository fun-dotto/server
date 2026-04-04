package service

import "github.com/fun-dotto/api-template/internal/domain"

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}
