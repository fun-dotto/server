package model

import "github.com/google/uuid"

type TimetableItem struct {
	Common

	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject   *Subject  `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DayOfWeek *string
	Period    *string
	// Deprecated: Use many-to-many relationship instead
	// Rooms     []Room `gorm:"many2many:timetable_item_rooms"`
	Rooms []TimetableItemRoom `gorm:"foreignKey:TimetableItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TimetableItemRoom struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	TimetableItemID string `gorm:"type:uuid;not null;index"`
	RoomID          string `gorm:"type:uuid;not null;index"`
	Room            *Room  `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
