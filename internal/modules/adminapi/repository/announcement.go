package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewAnnouncementRepository(db *gorm.DB) *Repository[model.Announcement] {
	return newRepository[model.Announcement](db, []string{"id"})
}
