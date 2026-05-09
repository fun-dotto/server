package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type TimetableItem struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	SubjectID string `gorm:"type:uuid;not null;index"`
	DayOfWeek *string
	Period    *string
	Subject   *Subject            `gorm:"foreignKey:SubjectID"`
	Rooms     []TimetableItemRoom `gorm:"foreignKey:TimetableItemID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TimetableItemRoom struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	TimetableItemID string `gorm:"type:uuid;not null;index"`
	RoomID          string `gorm:"type:uuid;not null;index"`
	Room            *Room  `gorm:"foreignKey:RoomID"`
}

func TimetableItemToDomain(m TimetableItem) domain.TimetableItem {
	var slot *domain.TimetableSlot
	if m.DayOfWeek != nil && m.Period != nil {
		slot = &domain.TimetableSlot{
			DayOfWeek: domain.DayOfWeek(*m.DayOfWeek),
			Period:    domain.Period(*m.Period),
		}
	}

	var subject domain.Subject
	if m.Subject != nil {
		subject = SubjectToDomain(*m.Subject)
	}

	rooms := make([]domain.Room, 0, len(m.Rooms))
	for _, r := range m.Rooms {
		if r.Room != nil {
			rooms = append(rooms, RoomToDomain(*r.Room))
		}
	}

	return domain.TimetableItem{
		ID:      m.ID,
		Subject: subject,
		Slot:    slot,
		Rooms:   rooms,
	}
}

func TimetableItemFromDomain(d domain.TimetableItem) TimetableItem {
	var dayOfWeek *string
	var period *string

	if d.Slot != nil {
		dow := string(d.Slot.DayOfWeek)
		p := string(d.Slot.Period)
		dayOfWeek = &dow
		period = &p
	}

	rooms := make([]TimetableItemRoom, len(d.Rooms))
	for i, room := range d.Rooms {
		rooms[i] = TimetableItemRoom{
			TimetableItemID: d.ID,
			RoomID:          room.ID,
		}
	}

	return TimetableItem{
		ID:        d.ID,
		SubjectID: d.Subject.ID,
		DayOfWeek: dayOfWeek,
		Period:    period,
		Rooms:     rooms,
	}
}
