package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *UserRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	var dbUsers []model.User
	if err := r.db.WithContext(ctx).Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, userToDomain(u))
	}

	return users, nil
}
