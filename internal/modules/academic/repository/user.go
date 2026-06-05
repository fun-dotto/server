package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID はユーザーを取得する。見つからない場合は found=false を返す。
func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, bool, error) {
	var record model.User
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}
	return userToDomain(record), true, nil
}

func userToDomain(m model.User) domain.User {
	u := domain.User{ID: m.ID}
	if m.Course != nil {
		c := domain.CourseType(*m.Course)
		u.Course = &c
	}
	if m.Grade != nil {
		g := domain.Grade(*m.Grade)
		u.Grade = &g
	}
	return u
}
