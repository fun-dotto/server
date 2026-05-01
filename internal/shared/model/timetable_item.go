package model

import "github.com/google/uuid"

type TimetableItem struct {
	Common

	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	DayOfWeek *string
	Period    *string
	Subject   *Subject `gorm:"foreignKey:SubjectID"`
	Rooms     []Room   `gorm:"many2many:timetable_item_rooms"`
}
