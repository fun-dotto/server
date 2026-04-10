package service

import (
	"github.com/fun-dotto/user-api/internal/domain"
)

type UserRepository interface {
	ListUsers() ([]domain.User, error)
	GetUserByID(id string) (domain.User, error)
	UpsertUser(user domain.User) (domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id string) (domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpsertUser(user domain.User) (domain.User, error) {
	return s.repo.UpsertUser(user)
}
