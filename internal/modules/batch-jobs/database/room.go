package database

import "github.com/fun-dotto/server/internal/modules/batch-jobs/domain"

type Room struct {
	ID   string `gorm:"type:uuid;primaryKey"`
	Name string `gorm:"not null"`
}

func (m *Room) ToDomain() domain.Room {
	return domain.Room{
		ID:   m.ID,
		Name: m.Name,
	}
}
