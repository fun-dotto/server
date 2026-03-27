package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type Room struct {
	ID    string `gorm:"type:uuid;primaryKey"`
	Name  string `gorm:"not null"`
	Floor string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func RoomToDomain(m Room) domain.Room {
	return domain.Room{
		ID:    m.ID,
		Name:  m.Name,
		Floor: domain.Floor(m.Floor),
	}
}

func RoomFromDomain(d domain.Room) Room {
	return Room{
		ID:    d.ID,
		Name:  d.Name,
		Floor: string(d.Floor),
	}
}
