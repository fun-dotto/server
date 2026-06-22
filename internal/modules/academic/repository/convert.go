package repository

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

// shared/model は ID を uuid.UUID で保持する一方、academic の domain 層は
// 文字列 ID を扱う。境界変換をこのファイルに集約する。

func parseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func parseUUIDs(ss []string) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(ss))
	for _, s := range ss {
		if id, err := uuid.Parse(s); err == nil {
			out = append(out, id)
		}
	}
	return out
}

// parseDomainDate は domain 層の YYYY-MM-DD 文字列を time.Time に変換する。
// パース不能な場合は zero value を返す（既存挙動を維持）。
func parseDomainDate(s string) time.Time {
	t, _ := time.Parse(dateLayout, s)
	return t
}

// 複合 PK を持つテーブル（CourseRegistration / FacultyRoom）は、API 側の
// 単一 ID 表現と整合させるため "a|b" 形式で文字列化する。

const compositeIDSep = "|"

// ErrInvalidCompositeID は複合 ID 文字列の形式が不正な場合に返される。
var ErrInvalidCompositeID = errors.New("invalid composite id")

func encodeCourseRegistrationID(userID string, subjectID uuid.UUID) string {
	return userID + compositeIDSep + subjectID.String()
}

func decodeCourseRegistrationID(id string) (userID string, subjectID uuid.UUID, err error) {
	parts := strings.SplitN(id, compositeIDSep, 2)
	if len(parts) != 2 {
		return "", uuid.Nil, ErrInvalidCompositeID
	}
	sid, perr := uuid.Parse(parts[1])
	if perr != nil {
		return "", uuid.Nil, ErrInvalidCompositeID
	}
	return parts[0], sid, nil
}

func encodeFacultyRoomID(facultyID, roomID uuid.UUID, year int) string {
	return facultyID.String() + compositeIDSep + roomID.String() + compositeIDSep + strconv.Itoa(year)
}

func decodeFacultyRoomID(id string) (facultyID, roomID uuid.UUID, year int, err error) {
	parts := strings.Split(id, compositeIDSep)
	if len(parts) != 3 {
		return uuid.Nil, uuid.Nil, 0, ErrInvalidCompositeID
	}
	fid, ferr := uuid.Parse(parts[0])
	if ferr != nil {
		return uuid.Nil, uuid.Nil, 0, ErrInvalidCompositeID
	}
	rid, rerr := uuid.Parse(parts[1])
	if rerr != nil {
		return uuid.Nil, uuid.Nil, 0, ErrInvalidCompositeID
	}
	y, yerr := strconv.Atoi(parts[2])
	if yerr != nil {
		return uuid.Nil, uuid.Nil, 0, ErrInvalidCompositeID
	}
	return fid, rid, y, nil
}

// --- entity converters ---

func facultyToDomain(m model.Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    m.ID.String(),
		Name:  m.Name,
		Email: m.Email,
	}
}

func facultyFromDomain(d domain.Faculty) model.Faculty {
	m := model.Faculty{
		Name:  d.Name,
		Email: d.Email,
	}
	m.ID = parseUUIDOrNil(d.ID)
	return m
}

func roomToDomain(m model.Room) domain.Room {
	return domain.Room{
		ID:    m.ID.String(),
		Name:  m.Name,
		Floor: domain.Floor(m.Floor),
	}
}

func roomFromDomain(d domain.Room) model.Room {
	m := model.Room{
		Name:  d.Name,
		Floor: string(d.Floor),
	}
	m.ID = parseUUIDOrNil(d.ID)
	return m
}

func syllabusToDomain(m model.Syllabus) domain.Syllabus {
	return domain.Syllabus{
		ID:                           m.ID,
		Name:                         m.Name,
		EnName:                       m.EnName,
		Grades:                       m.Grades,
		Credit:                       m.Credit,
		FacultyNames:                 m.FacultyNames,
		PracticalHomeFacultyCategory: m.PracticalHomeFacultyCategory,
		MultiplePersonTeachingForm:   m.MultiplePersonTeachingForm,
		TeachingForm:                 m.TeachingForm,
		Summary:                      m.Summary,
		LearningOutcomes:             m.LearningOutcomes,
		Assignments:                  m.Assignments,
		EvaluationMethod:             m.EvaluationMethod,
		Textbooks:                    m.Textbooks,
		ReferenceBooks:               m.ReferenceBooks,
		Prerequisites:                m.Prerequisites,
		PreLearning:                  m.PreLearning,
		PostLearning:                 m.PostLearning,
		Notes:                        m.Notes,
		Keywords:                     m.Keywords,
		TargetCourses:                m.TargetCourses,
		TargetAreas:                  m.TargetAreas,
		Classifications:              m.Classifications,
		TeachingLanguage:             m.TeachingLanguage,
		ContentsAndSchedule:          m.ContentsAndSchedule,
		TeachingAndExamForm:          m.TeachingAndExamForm,
		DsopSubject:                  m.DsopSubject,
	}
}

func subjectToDomain(m model.Subject) domain.Subject {
	faculties := make([]domain.SubjectFaculty, len(m.Faculties))
	for i, f := range m.Faculties {
		var faculty domain.Faculty
		if f.Faculty != nil {
			faculty = facultyToDomain(*f.Faculty)
		} else {
			faculty = domain.Faculty{ID: f.FacultyID.String()}
		}
		faculties[i] = domain.SubjectFaculty{
			Faculty:   faculty,
			IsPrimary: f.IsPrimary,
		}
	}

	eligible := make([]domain.SubjectTargetClass, len(m.EligibleAttributes))
	for i, e := range m.EligibleAttributes {
		tc := domain.SubjectTargetClass{
			Grade: domain.Grade(e.Grade),
		}
		if e.Class != nil {
			c := domain.Class(*e.Class)
			tc.Class = &c
		}
		eligible[i] = tc
	}

	requirements := make([]domain.SubjectRequirement, len(m.Requirements))
	for i, r := range m.Requirements {
		requirements[i] = domain.SubjectRequirement{
			Course:          domain.CourseType(r.Course),
			RequirementType: domain.SubjectRequirementType(r.RequirementType),
		}
	}

	return domain.Subject{
		ID:                      m.ID.String(),
		Name:                    m.Name,
		Faculties:               faculties,
		Year:                    m.Year,
		Semester:                domain.CourseSemester(m.Semester),
		Credit:                  m.Credit,
		Classification:          domain.SubjectClassification(m.Classification),
		CulturalSubjectCategory: domain.CulturalSubjectCategory(m.CulturalSubjectCategory),
		EligibleAttributes:      eligible,
		Requirements:            requirements,
		SyllabusID:              m.SyllabusID,
	}
}

func cancelledClassToDomain(m model.CancelledClass) domain.CancelledClass {
	d := domain.CancelledClass{
		ID:     m.ID.String(),
		Date:   m.Date.Format(dateLayout),
		Period: domain.Period(m.Period),
	}
	if m.Comment != nil {
		d.Comment = *m.Comment
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(*m.Subject)
	}
	return d
}

func cancelledClassFromDomain(d domain.CancelledClass) model.CancelledClass {
	var comment *string
	if d.Comment != "" {
		comment = &d.Comment
	}
	m := model.CancelledClass{
		SubjectID: parseUUIDOrNil(d.Subject.ID),
		Date:      parseDomainDate(d.Date),
		Period:    string(d.Period),
		Comment:   comment,
	}
	if d.ID != "" {
		m.ID = parseUUIDOrNil(d.ID)
	}
	return m
}

func makeupClassToDomain(m model.MakeupClass) domain.MakeupClass {
	d := domain.MakeupClass{
		ID:     m.ID.String(),
		Date:   m.Date.Format(dateLayout),
		Period: domain.Period(m.Period),
	}
	if m.Comment != nil {
		d.Comment = *m.Comment
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(*m.Subject)
	}
	return d
}

func makeupClassFromDomain(d domain.MakeupClass) model.MakeupClass {
	var comment *string
	if d.Comment != "" {
		comment = &d.Comment
	}
	m := model.MakeupClass{
		SubjectID: parseUUIDOrNil(d.Subject.ID),
		Date:      parseDomainDate(d.Date),
		Period:    string(d.Period),
		Comment:   comment,
	}
	if d.ID != "" {
		m.ID = parseUUIDOrNil(d.ID)
	}
	return m
}

func roomChangeToDomain(m model.RoomChange) domain.RoomChange {
	d := domain.RoomChange{
		ID:     m.ID.String(),
		Date:   m.Date.Format(dateLayout),
		Period: domain.Period(m.Period),
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(*m.Subject)
	}
	if m.OriginalRoom != nil {
		d.OriginalRoom = roomToDomain(*m.OriginalRoom)
	}
	if m.NewRoom != nil {
		d.NewRoom = roomToDomain(*m.NewRoom)
	}
	return d
}

func roomChangeFromDomain(d domain.RoomChange) model.RoomChange {
	m := model.RoomChange{
		SubjectID:      parseUUIDOrNil(d.Subject.ID),
		Date:           parseDomainDate(d.Date),
		Period:         string(d.Period),
		OriginalRoomID: parseUUIDOrNil(d.OriginalRoom.ID),
		NewRoomID:      parseUUIDOrNil(d.NewRoom.ID),
	}
	if d.ID != "" {
		m.ID = parseUUIDOrNil(d.ID)
	}
	return m
}

func courseRegistrationToDomain(m model.CourseRegistration) domain.CourseRegistration {
	d := domain.CourseRegistration{
		ID:     encodeCourseRegistrationID(m.UserID, m.SubjectID),
		UserID: m.UserID,
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(*m.Subject)
	} else {
		d.Subject = domain.Subject{ID: m.SubjectID.String()}
	}
	return d
}

func courseRegistrationFromDomain(d domain.CourseRegistration) model.CourseRegistration {
	return model.CourseRegistration{
		UserID:    d.UserID,
		SubjectID: parseUUIDOrNil(d.Subject.ID),
	}
}

func facultyRoomToDomain(m model.FacultyRoom) domain.FacultyRoom {
	d := domain.FacultyRoom{
		ID:   encodeFacultyRoomID(m.FacultyID, m.RoomID, m.Year),
		Year: m.Year,
	}
	if m.Faculty != nil {
		d.Faculty = facultyToDomain(*m.Faculty)
	} else {
		d.Faculty = domain.Faculty{ID: m.FacultyID.String()}
	}
	if m.Room != nil {
		d.Room = roomToDomain(*m.Room)
	} else {
		d.Room = domain.Room{ID: m.RoomID.String()}
	}
	return d
}

func facultyRoomFromDomain(d domain.FacultyRoom) model.FacultyRoom {
	return model.FacultyRoom{
		FacultyID: parseUUIDOrNil(d.Faculty.ID),
		RoomID:    parseUUIDOrNil(d.Room.ID),
		Year:      d.Year,
	}
}

func timetableItemToDomain(m model.TimetableItem) domain.TimetableItem {
	var slot *domain.TimetableSlot
	if m.DayOfWeek != nil && m.Period != nil {
		slot = &domain.TimetableSlot{
			DayOfWeek: domain.DayOfWeek(*m.DayOfWeek),
			Period:    domain.Period(*m.Period),
		}
	}

	var subject domain.Subject
	if m.Subject != nil {
		subject = subjectToDomain(*m.Subject)
	}

	rooms := make([]domain.Room, 0, len(m.Rooms))
	for _, r := range m.Rooms {
		if r.Room != nil {
			rooms = append(rooms, roomToDomain(*r.Room))
		}
	}

	return domain.TimetableItem{
		ID:      m.ID.String(),
		Subject: subject,
		Slot:    slot,
		Rooms:   rooms,
	}
}

func timetableItemFromDomain(d domain.TimetableItem) model.TimetableItem {
	var dayOfWeek *string
	var period *string
	if d.Slot != nil {
		dow := string(d.Slot.DayOfWeek)
		p := string(d.Slot.Period)
		dayOfWeek = &dow
		period = &p
	}

	rooms := make([]model.TimetableItemRoom, len(d.Rooms))
	for i, room := range d.Rooms {
		rm := roomFromDomain(room)
		rooms[i] = model.TimetableItemRoom{
			RoomID: rm.ID.String(),
			Room:   &rm,
		}
	}

	m := model.TimetableItem{
		SubjectID: parseUUIDOrNil(d.Subject.ID),
		DayOfWeek: dayOfWeek,
		Period:    period,
		Rooms:     rooms,
	}
	if d.ID != "" {
		m.ID = parseUUIDOrNil(d.ID)
	}
	return m
}
