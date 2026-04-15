package service

import (
	"context"

	"github.com/fun-dotto/user-api/internal/domain"
)

type FCMTokenRepository interface {
	ListFCMTokens(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error)
	UpsertFCMToken(ctx context.Context, token domain.FCMToken) (domain.FCMToken, error)
}

type FCMTokenService struct {
	repo FCMTokenRepository
}

func NewFCMTokenService(repo FCMTokenRepository) *FCMTokenService {
	return &FCMTokenService{repo: repo}
}

func (s *FCMTokenService) ListFCMTokens(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
	return s.repo.ListFCMTokens(ctx, filter)
}

func (s *FCMTokenService) UpsertFCMToken(ctx context.Context, token domain.FCMToken) (domain.FCMToken, error) {
	return s.repo.UpsertFCMToken(ctx, token)
}
