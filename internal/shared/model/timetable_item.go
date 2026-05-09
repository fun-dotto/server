package model

import "github.com/google/uuid"

type TimetableItem struct {
	Common

	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID"`
	DayOfWeek *string
	Period    *string
	Rooms     []Room `gorm:"many2many:timetable_item_rooms"`
}
