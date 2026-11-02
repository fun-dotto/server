package model

import (
	"time"
	"uuid"
)

type CourseRegistration struct {
	// Deprecated: Id field is no longer used
	Id        uuid.UUID `gorm:"type:uuid"`
	UserID    string    `gorm:"primaryKey"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
