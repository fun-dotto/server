package model

import (
	"time"

	"github.com/google/uuid"
)

type CourseRegistration struct {
	UserID    string    `gorm:"primaryKey"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE"`
	SubjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
