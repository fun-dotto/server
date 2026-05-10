package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type FCMTokenRepository struct {
	db *gorm.DB
}

func NewFCMTokenRepository(db *gorm.DB) *FCMTokenRepository {
	return &FCMTokenRepository{db: db}
}

func (r *FCMTokenRepository) ListFCMTokens(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
	query := r.db.WithContext(ctx).Model(&model.FCMToken{})

	if len(filter.UserIDs) > 0 {
		query = query.Where("user_id IN ?", filter.UserIDs)
	}
	if len(filter.Tokens) > 0 {
		query = query.Where("token IN ?", filter.Tokens)
	}
	if filter.UpdatedAtFrom != nil {
		query = query.Where("updated_at >= ?", *filter.UpdatedAtFrom)
	}
	if filter.UpdatedAtTo != nil {
		query = query.Where("updated_at <= ?", *filter.UpdatedAtTo)
	}

	var dbTokens []model.FCMToken
	if err := query.Order("updated_at DESC").Find(&dbTokens).Error; err != nil {
		return nil, err
	}

	tokens := make([]domain.FCMToken, 0, len(dbTokens))
	for i := range dbTokens {
		tokens = append(tokens, fcmTokenToDomain(&dbTokens[i]))
	}

	return tokens, nil
}
