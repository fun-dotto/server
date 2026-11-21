package model

import (
	"time"
	"uuid"
)

type FacultyRoom struct {
	// Deprecated: 列を DROP するまでモデル差分チェックを通すために残している。読み書きしない
	Id        uuid.UUID `gorm:"type:uuid;->:false"`
	FacultyID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Faculty   *Faculty  `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoomID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Room      *Room     `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Year      int       `gorm:"not null;primaryKey"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
