package repository

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"gorm.io/gorm"
)

type DayOfWeekTimetableSlotRepository struct {
	db *gorm.DB
}

func NewDayOfWeekTimetableSlotRepository(db *gorm.DB) *DayOfWeekTimetableSlotRepository {
	return &DayOfWeekTimetableSlotRepository{db: db}
}

func (r *DayOfWeekTimetableSlotRepository) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	var records []database.DayOfWeekTimetableSlot
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.DayOfWeekTimetableSlot, len(records))
	for i, rec := range records {
		results[i] = database.DayOfWeekTimetableSlotToDomain(rec)
	}
	return results, nil
}

func (r *DayOfWeekTimetableSlotRepository) GetByID(ctx context.Context, id string) (domain.DayOfWeekTimetableSlot, error) {
	var record database.DayOfWeekTimetableSlot
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		return domain.DayOfWeekTimetableSlot{}, err
	}
	return database.DayOfWeekTimetableSlotToDomain(record), nil
}

func (r *DayOfWeekTimetableSlotRepository) Create(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error) {
	record := database.DayOfWeekTimetableSlotFromDomain(slot)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.DayOfWeekTimetableSlot{}, err
	}
	return database.DayOfWeekTimetableSlotToDomain(record), nil
}

func (r *DayOfWeekTimetableSlotRepository) Update(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error) {
	record := database.DayOfWeekTimetableSlotFromDomain(slot)
	if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
		return domain.DayOfWeekTimetableSlot{}, err
	}
	return database.DayOfWeekTimetableSlotToDomain(record), nil
}

func (r *DayOfWeekTimetableSlotRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.DayOfWeekTimetableSlot{}, "id = ?", id).Error
}
