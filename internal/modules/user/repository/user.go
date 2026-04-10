package repository

import (
	"errors"

	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(id string) (domain.User, error) {
	var user database.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) UpsertUser(user domain.User) (domain.User, error) {
	dbUser := database.UserFromDomain(user)
	if err := r.db.Save(&dbUser).Error; err != nil {
		return domain.User{}, err
	}
	return dbUser.ToDomain(), nil
}
