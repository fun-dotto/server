package model

import (
	"time"

	"github.com/google/uuid"
)

type RoomChange struct {
	Common

	SubjectID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject        *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`
	Date           time.Time `gorm:"type:date;not null;index"`
	Period         string    `gorm:"not null"`
	OriginalRoomID uuid.UUID `gorm:"type:uuid;not null"`
	OriginalRoom   *Room     `gorm:"foreignKey:OriginalRoomID;constraint:OnUpdate:CASCADE"`
	NewRoomID      uuid.UUID `gorm:"type:uuid;not null"`
	NewRoom        *Room     `gorm:"foreignKey:NewRoomID;constraint:OnUpdate:CASCADE"`
}
