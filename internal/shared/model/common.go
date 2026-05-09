package model

import (
	"time"

	"github.com/google/uuid"
)

type Common struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
