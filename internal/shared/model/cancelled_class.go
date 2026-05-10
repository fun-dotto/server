package model

import "github.com/google/uuid"

type CancelledClass struct {
	Common

	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`
	Date      string    `gorm:"type:date;not null;index"`
	Period    string    `gorm:"not null"`
	Comment   *string
}
