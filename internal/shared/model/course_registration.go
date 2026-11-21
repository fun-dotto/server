package model

import (
	"time"
	"uuid"
)

type CourseRegistration struct {
	// Deprecated: 列を DROP するまでモデル差分チェックを通すために残している。読み書きしない
	Id        uuid.UUID `gorm:"type:uuid;->:false"`
	UserID    string    `gorm:"primaryKey"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
