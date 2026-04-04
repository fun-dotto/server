package service

import "github.com/fun-dotto/api-template/internal/domain"

type FCMTokenRepository interface {
	ListFCMTokens(filter domain.FCMTokenListFilter) ([]domain.FCMToken, error)
	UpsertFCMToken(token domain.FCMToken) (domain.FCMToken, error)
}

type FCMTokenService struct {
	repo FCMTokenRepository
}

func NewFCMTokenService(repo FCMTokenRepository) *FCMTokenService {
	return &FCMTokenService{repo: repo}
}

func (s *FCMTokenService) ListFCMTokens(filter domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
	return s.repo.ListFCMTokens(filter)
}

func (s *FCMTokenService) UpsertFCMToken(token domain.FCMToken) (domain.FCMToken, error) {
	return s.repo.UpsertFCMToken(token)
}
