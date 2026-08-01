package service

import "github.com/fun-dotto/server/internal/modules/app/domain"

type UserRepository interface {
	GetUser(id string) (*domain.User, error)
	UpsertUser(id string, req domain.UserRequest) (*domain.User, error)
	UpsertFCMToken(userID string, req domain.FCMTokenRequest) (*domain.FCMToken, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUser(id string) (*domain.User, error) {
	return s.repository.GetUser(id)
}

func (s *UserService) UpsertUser(id string, req domain.UserRequest) (*domain.User, error) {
	return s.repository.UpsertUser(id, req)
}

func (s *UserService) UpsertFCMToken(userID string, req domain.FCMTokenRequest) (*domain.FCMToken, error) {
	return s.repository.UpsertFCMToken(userID, req)
}
