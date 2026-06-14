package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	for _, t := range dbTokens {
		tokens = append(tokens, fcmTokenToDomain(t))
	}

	return tokens, nil
}

func (r *FCMTokenRepository) UpsertFCMToken(ctx context.Context, token domain.FCMToken) (domain.FCMToken, error) {
	dbToken := fcmTokenFromDomain(token)

	if err := r.db.WithContext(ctx).Omit("User").Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "token"}},
		DoUpdates: clause.Assignments(map[string]any{
			"user_id":    dbToken.UserID,
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&dbToken).Error; err != nil {
		return domain.FCMToken{}, err
	}

	if err := r.db.WithContext(ctx).First(&dbToken, "token = ?", dbToken.Token).Error; err != nil {
		return domain.FCMToken{}, err
	}

	return fcmTokenToDomain(dbToken), nil
}
