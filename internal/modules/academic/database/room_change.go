package database

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type RoomChange struct {
	ID             string    `gorm:"type:uuid;primaryKey"`
	SubjectID      string    `gorm:"type:uuid;not null;index"`
	Subject        *Subject  `gorm:"foreignKey:SubjectID"`
	Date           string    `gorm:"type:date;not null;index"`
	Period         string    `gorm:"not null"`
	OriginalRoomID string    `gorm:"type:uuid;not null"`
	OriginalRoom   *Room     `gorm:"foreignKey:OriginalRoomID"`
	NewRoomID      string    `gorm:"type:uuid;not null"`
	NewRoom        *Room     `gorm:"foreignKey:NewRoomID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func RoomChangeToDomain(m RoomChange) domain.RoomChange {
	date := m.Date
	if t, err := time.Parse("2006-01-02T15:04:05Z", m.Date); err == nil {
		date = t.Format("2006-01-02")
	} else if t, err := time.Parse(time.RFC3339, m.Date); err == nil {
		date = t.Format("2006-01-02")
	}
	d := domain.RoomChange{
		ID:     m.ID,
		Date:   date,
		Period: domain.Period(m.Period),
	}
	if m.Subject != nil {
		d.Subject = SubjectToDomain(*m.Subject)
	}
	if m.OriginalRoom != nil {
		d.OriginalRoom = RoomToDomain(*m.OriginalRoom)
	}
	if m.NewRoom != nil {
		d.NewRoom = RoomToDomain(*m.NewRoom)
	}
	return d
}

func RoomChangeFromDomain(d domain.RoomChange) RoomChange {
	return RoomChange{
		ID:             d.ID,
		SubjectID:      d.Subject.ID,
		Date:           d.Date,
		Period:         string(d.Period),
		OriginalRoomID: d.OriginalRoom.ID,
		NewRoomID:      d.NewRoom.ID,
	}
}
