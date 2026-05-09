package database

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type FacultyRoom struct {
	ID        string   `gorm:"type:uuid;primaryKey"`
	FacultyID string   `gorm:"type:uuid;not null;uniqueIndex:idx_faculty_rooms_faculty_year,priority:1"`
	Faculty   *Faculty `gorm:"foreignKey:FacultyID"`
	RoomID    string   `gorm:"type:uuid;not null;uniqueIndex:idx_faculty_rooms_room_year,priority:1"`
	Room      *Room    `gorm:"foreignKey:RoomID"`
	Year      int      `gorm:"not null;uniqueIndex:idx_faculty_rooms_faculty_year,priority:2;uniqueIndex:idx_faculty_rooms_room_year,priority:2"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func FacultyRoomToDomain(m FacultyRoom) domain.FacultyRoom {
	d := domain.FacultyRoom{
		ID:   m.ID,
		Year: m.Year,
	}
	if m.Faculty != nil {
		d.Faculty = FacultyToDomain(*m.Faculty)
	}
	if m.Room != nil {
		d.Room = RoomToDomain(*m.Room)
	}
	return d
}

func FacultyRoomFromDomain(d domain.FacultyRoom) FacultyRoom {
	return FacultyRoom{
		ID:        d.ID,
		FacultyID: d.Faculty.ID,
		RoomID:    d.Room.ID,
		Year:      d.Year,
	}
}
