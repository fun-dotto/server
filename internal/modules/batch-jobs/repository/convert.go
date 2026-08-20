package repository

import (
	"time"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

// parseUUIDOrNil は文字列を uuid.UUID に変換する。パース不能な場合は uuid.Nil を返す。
func parseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func cancelledClassFromDomain(d academicdomain.CancelledClass) model.CancelledClass {
	var comment *string
	if d.Comment != "" {
		comment = &d.Comment
	}
	return model.CancelledClass{
		SubjectID: parseUUIDOrNil(d.Subject.ID),
		Date:      parseDomainDate(d.Date),
		Period:    string(d.Period),
		Comment:   comment,
	}
}

func makeupClassFromDomain(d academicdomain.MakeupClass) model.MakeupClass {
	var comment *string
	if d.Comment != "" {
		comment = &d.Comment
	}
	return model.MakeupClass{
		SubjectID: parseUUIDOrNil(d.Subject.ID),
		Date:      parseDomainDate(d.Date),
		Period:    string(d.Period),
		Comment:   comment,
	}
}

func roomChangeFromDomain(d academicdomain.RoomChange) model.RoomChange {
	return model.RoomChange{
		SubjectID:      parseUUIDOrNil(d.Subject.ID),
		Date:           parseDomainDate(d.Date),
		Period:         string(d.Period),
		OriginalRoomID: parseUUIDOrNil(d.OriginalRoom.ID),
		NewRoomID:      parseUUIDOrNil(d.NewRoom.ID),
	}
}

func cancelledClassToDomain(m model.CancelledClass) academicdomain.CancelledClass {
	d := academicdomain.CancelledClass{
		ID:     m.ID.String(),
		Date:   m.Date.Format(dateLayout),
		Period: academicdomain.Period(m.Period),
	}
	if m.Comment != nil {
		d.Comment = *m.Comment
	}
	d.Subject = academicdomain.Subject{ID: m.SubjectID.String()}
	return d
}

func makeupClassToDomain(m model.MakeupClass) academicdomain.MakeupClass {
	d := academicdomain.MakeupClass{
		ID:     m.ID.String(),
		Date:   m.Date.Format(dateLayout),
		Period: academicdomain.Period(m.Period),
	}
	if m.Comment != nil {
		d.Comment = *m.Comment
	}
	d.Subject = academicdomain.Subject{ID: m.SubjectID.String()}
	return d
}

func roomChangeToDomain(m model.RoomChange) academicdomain.RoomChange {
	return academicdomain.RoomChange{
		ID:           m.ID.String(),
		Date:         m.Date.Format(dateLayout),
		Period:       academicdomain.Period(m.Period),
		Subject:      academicdomain.Subject{ID: m.SubjectID.String()},
		OriginalRoom: academicdomain.Room{ID: m.OriginalRoomID.String()},
		NewRoom:      academicdomain.Room{ID: m.NewRoomID.String()},
	}
}

func facultyRoomFromDomain(facultyID, roomID string, year int) model.FacultyRoom {
	return model.FacultyRoom{
		FacultyID: parseUUIDOrNil(facultyID),
		RoomID:    parseUUIDOrNil(roomID),
		Year:      year,
	}
}

// parseDomainDate は domain 層の YYYY-MM-DD 文字列を time.Time に変換する。
// パース不能な場合は zero value を返す。
func parseDomainDate(s string) time.Time {
	t, _ := time.Parse(dateLayout, s)
	return t
}
