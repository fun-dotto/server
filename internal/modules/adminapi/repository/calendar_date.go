package repository

import (
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

func NewCalendarDateRepository(db *gorm.DB) *Repository[model.CalendarDate] {
	return newRepository[model.CalendarDate](db, []string{"service_id", "date"})
}
