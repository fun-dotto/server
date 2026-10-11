package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewStopTimeRepository(db *gorm.DB) *Repository[model.StopTime] {
	return newRepository[model.StopTime](db, []string{"trip_id", "stop_sequence"})
}
