package model

import "time"

type FCMToken struct {
	Token  string `gorm:"primaryKey"`
	UserID string `gorm:"not null;index"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
