package repository

import (
	"github.com/fun-dotto/api-template/internal/database"
	"github.com/fun-dotto/api-template/internal/domain"
)

func (r *UserRepository) ListUsers() ([]domain.User, error) {
	var dbUsers []database.User
	if err := r.db.Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, u.ToDomain())
	}

	return users, nil
}
