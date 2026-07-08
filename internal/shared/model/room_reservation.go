package model

import (
	"time"

	"github.com/google/uuid"
)

type RoomReservation struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Title     string    `gorm:"not null"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null"`
	Room      *Room     `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE"`
	StartTime time.Time `gorm:"not null;index"`
	EndTime   time.Time `gorm:"not null;index"`
}
