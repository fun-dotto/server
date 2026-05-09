package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTargetUser struct {
	NotificationID uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Notification   Notification `gorm:"foreignKey:NotificationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID         string       `gorm:"primaryKey"`
	User           User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	NotifiedAt     *time.Time   `gorm:"index"`
}
