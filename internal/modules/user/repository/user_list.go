package repository

import (
	"context"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (r *UserRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	var dbUsers []database.User
	if err := r.db.WithContext(ctx).Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, u.ToDomain())
	}

	return users, nil
}
