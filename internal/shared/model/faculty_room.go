package model

import (
	"time"

	"github.com/google/uuid"
)

type FacultyRoom struct {
	FacultyID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Faculty   *Faculty  `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE"`
	RoomID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Room      *Room     `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE"`
	Year      int       `gorm:"not null;primaryKey"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
