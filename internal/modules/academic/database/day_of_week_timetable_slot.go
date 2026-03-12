package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type DayOfWeekTimetableSlot struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	DayOfWeek     string `gorm:"type:text;not null"`
	TimetableSlot string `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func DayOfWeekTimetableSlotToDomain(m DayOfWeekTimetableSlot) domain.DayOfWeekTimetableSlot {
	return domain.DayOfWeekTimetableSlot{
		ID:            m.ID,
		DayOfWeek:     domain.DayOfWeek(m.DayOfWeek),
		TimetableSlot: domain.TimetableSlot(m.TimetableSlot),
	}
}

func DayOfWeekTimetableSlotFromDomain(d domain.DayOfWeekTimetableSlot) DayOfWeekTimetableSlot {
	return DayOfWeekTimetableSlot{
		ID:            d.ID,
		DayOfWeek:     string(d.DayOfWeek),
		TimetableSlot: string(d.TimetableSlot),
	}
}
