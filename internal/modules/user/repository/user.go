package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return userToDomain(user), nil
}

func (r *UserRepository) UpsertUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser := userFromDomain(user)
	if err := r.db.WithContext(ctx).Save(&dbUser).Error; err != nil {
		return domain.User{}, err
	}
	return userToDomain(dbUser), nil
}
