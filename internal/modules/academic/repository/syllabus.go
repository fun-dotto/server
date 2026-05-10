package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type SyllabusRepository struct {
	db *gorm.DB
}

func NewSyllabusRepository(db *gorm.DB) *SyllabusRepository {
	return &SyllabusRepository{db: db}
}

func (r *SyllabusRepository) GetByID(ctx context.Context, id string) (domain.Syllabus, error) {
	var record model.Syllabus
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return domain.Syllabus{}, err
	}
	return syllabusToDomain(record), nil
}
