package model

import (
	"time"
	"uuid"
)

type RoomReservation struct {
	Common

	Title     string    `gorm:"not null"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Room      *Room     `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE"`
	StartTime time.Time `gorm:"not null;index"`
	EndTime   time.Time `gorm:"not null;index"`
}
