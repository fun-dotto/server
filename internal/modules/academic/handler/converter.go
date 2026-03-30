package handler

import (
	"strings"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

// Subject

func subjectToAPI(d domain.Subject) api.Subject {
	eligible := make([]api.SubjectTargetClass, len(d.EligibleAttributes))
	for i, e := range d.EligibleAttributes {
		tc := api.SubjectTargetClass{
			Grade: api.DottoFoundationV1Grade(e.Grade),
		}
		if e.Class != nil {
			c := api.DottoFoundationV1Class(*e.Class)
			tc.Class = &c
		}
		eligible[i] = tc
	}

	requirements := make([]api.SubjectRequirement, len(d.Requirements))
	for i, r := range d.Requirements {
		requirements[i] = api.SubjectRequirement{
			Course:          api.DottoFoundationV1Course(r.Course),
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}

	faculties := make([]api.SubjectFaculty, len(d.Faculties))
	for i, f := range d.Faculties {
		faculties[i] = api.SubjectFaculty{
			Faculty:   facultyToAPI(f.Faculty),
			IsPrimary: f.IsPrimary,
		}
	}

	return api.Subject{
		Id:                 d.ID,
		Name:               d.Name,
		Faculties:          faculties,
		Year:               d.Year,
		Semester:           api.DottoFoundationV1CourseSemester(d.Semester),
		Credit:             d.Credit,
		EligibleAttributes: eligible,
		Requirements:       requirements,
	}
}

func subjectToSummaryAPI(d domain.Subject) api.SubjectSummary {
	faculties := make([]api.SubjectFaculty, len(d.Faculties))
	for i, f := range d.Faculties {
		faculties[i] = api.SubjectFaculty{
			Faculty:   facultyToAPI(f.Faculty),
			IsPrimary: f.IsPrimary,
		}
	}

	return api.SubjectSummary{
		Id:        d.ID,
		Name:      d.Name,
		Faculties: faculties,
	}
}

func subjectsToSummaryAPI(ds []domain.Subject) []api.SubjectSummary {
	result := make([]api.SubjectSummary, len(ds))
	for i, d := range ds {
		result[i] = subjectToSummaryAPI(d)
	}
	return result
}

func buildSubjectListFilter(params api.SubjectsV1ListParams) domain.SubjectListFilter {
	filter := domain.SubjectListFilter{}

	if params.Ids != nil {
		filter.IDs = *params.Ids
	}
	if params.Q != nil {
		filter.Q = params.Q
	}
	if params.Grades != nil {
		grades := make([]domain.Grade, len(*params.Grades))
		for i, g := range *params.Grades {
			grades[i] = domain.Grade(g)
		}
		filter.Grade = grades
	}
	if params.Courses != nil {
		courses := make([]domain.CourseType, len(*params.Courses))
		for i, c := range *params.Courses {
			courses[i] = domain.CourseType(c)
		}
		filter.Courses = courses
	}
	if params.Classes != nil {
		classes := make([]domain.Class, len(*params.Classes))
		for i, c := range *params.Classes {
			classes[i] = domain.Class(c)
		}
		filter.Class = classes
	}
	if params.Classifications != nil {
		classifications := make([]domain.SubjectClassification, len(*params.Classifications))
		for i, c := range *params.Classifications {
			classifications[i] = domain.SubjectClassification(c)
		}
		filter.Classification = classifications
	}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		// デフォルトで今年度を設定
		// TODO: このデフォルト値設定のロジックは service 層に移すべき。
		// また、日本の大学の年度は4月始まりのため、1〜3月は前年度を返す必要がある。
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	if params.Semesters != nil {
		semesters := make([]domain.CourseSemester, len(*params.Semesters))
		for i, s := range *params.Semesters {
			semesters[i] = domain.CourseSemester(s)
		}
		filter.Semester = semesters
	}
	if params.RequirementTypes != nil {
		reqTypes := make([]domain.SubjectRequirementType, len(*params.RequirementTypes))
		for i, r := range *params.RequirementTypes {
			reqTypes[i] = domain.SubjectRequirementType(r)
		}
		filter.RequirementType = reqTypes
	}
	if params.CulturalSubjectCategories != nil {
		cats := make([]domain.CulturalSubjectCategory, len(*params.CulturalSubjectCategories))
		for i, c := range *params.CulturalSubjectCategories {
			cats[i] = domain.CulturalSubjectCategory(c)
		}
		filter.CulturalSubjectCategory = cats
	}

	return filter
}

// Syllabus

func syllabusToAPI(d domain.Syllabus) api.Syllabus {
	return api.Syllabus{
		Id:                           d.ID,
		Name:                         d.Name,
		EnName:                       d.EnName,
		Grades:                       d.Grades,
		Credit:                       d.Credit,
		FacultyNames:                 d.FacultyNames,
		PracticalHomeFacultyCategory: d.PracticalHomeFacultyCategory,
		MultiplePersonTeachingForm:   d.MultiplePersonTeachingForm,
		TeachingForm:                 d.TeachingForm,
		Summary:                      d.Summary,
		LearningOutcomes:             d.LearningOutcomes,
		Assignments:                  d.Assignments,
		EvaluationMethod:             d.EvaluationMethod,
		Textbooks:                    d.Textbooks,
		ReferenceBooks:               d.ReferenceBooks,
		Prerequisites:                d.Prerequisites,
		PreLearning:                  d.PreLearning,
		PostLearning:                 d.PostLearning,
		Notes:                        d.Notes,
		Keywords:                     d.Keywords,
		TargetCourses:                d.TargetCourses,
		TargetAreas:                  d.TargetAreas,
		Classifications:              d.Classifications,
		TeachingLanguage:             d.TeachingLanguage,
		ContentsAndSchedule:          d.ContentsAndSchedule,
		TeachingAndExamForm:          d.TeachingAndExamForm,
		DsopSubject:                  d.DsopSubject,
	}
}

func facultyToAPI(faculty domain.Faculty) api.Faculty {
	return api.Faculty{
		Id:    faculty.ID,
		Name:  faculty.Name,
		Email: faculty.Email,
	}
}

func facultiesToAPI(faculties []domain.Faculty) []api.Faculty {
	result := make([]api.Faculty, len(faculties))
	for i, faculty := range faculties {
		result[i] = facultyToAPI(faculty)
	}
	return result
}

func toDomainFacultyFromRequest(id string, req api.FacultyRequest) domain.Faculty {
	return domain.Faculty{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
}

// Room

func buildRoomListFilter(params api.RoomsV1ListParams) domain.RoomListFilter {
	filter := domain.RoomListFilter{}
	if params.Ids != nil {
		filter.IDs = *params.Ids
	}
	if params.Floor != nil {
		floors := make([]domain.Floor, len(*params.Floor))
		for i, f := range *params.Floor {
			floors[i] = domain.Floor(f)
		}
		filter.Floors = floors
	}
	return filter
}

func roomToAPI(room domain.Room) api.Room {
	return api.Room{
		Id:    room.ID,
		Name:  room.Name,
		Floor: api.DottoFoundationV1Floor(room.Floor),
	}
}

func roomsToAPI(rooms []domain.Room) []api.Room {
	result := make([]api.Room, len(rooms))
	for i, room := range rooms {
		result[i] = roomToAPI(room)
	}
	return result
}

func toDomainRoomFromRequest(id string, req api.RoomRequest) domain.Room {
	return domain.Room{
		ID:    id,
		Name:  req.Name,
		Floor: domain.Floor(req.Floor),
	}
}

// TimetableItem

func buildTimetableItemListFilter(params api.TimetableItemsV1ListParams) domain.TimetableItemListFilter {
	filter := domain.TimetableItemListFilter{}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		// TODO: このデフォルト値設定のロジックは service 層に移すべき。
		// また、日本の大学の年度は4月始まりのため、1〜3月は前年度を返す必要がある。
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	semesters := make([]domain.CourseSemester, len(params.Semesters))
	for i, s := range params.Semesters {
		semesters[i] = domain.CourseSemester(s)
	}
	filter.Semesters = semesters
	return filter
}

func timetableSlotToAPI(slot domain.TimetableSlot) api.DottoFoundationV1TimetableSlot {
	return api.DottoFoundationV1TimetableSlot{
		DayOfWeek: api.DottoFoundationV1DayOfWeek(slot.DayOfWeek),
		Period:    api.DottoFoundationV1Period(slot.Period),
	}
}

func timetableItemToAPI(d domain.TimetableItem) api.TimetableItem {
	var slot *api.DottoFoundationV1TimetableSlot
	if d.Slot != nil &&
		strings.TrimSpace(string(d.Slot.DayOfWeek)) != "" &&
		strings.TrimSpace(string(d.Slot.Period)) != "" {
		s := timetableSlotToAPI(*d.Slot)
		slot = &s
	}

	rooms := make([]api.Room, len(d.Rooms))
	for i, r := range d.Rooms {
		rooms[i] = roomToAPI(r)
	}

	return api.TimetableItem{
		Id:      d.ID,
		Subject: subjectToSummaryAPI(d.Subject),
		Slot:    slot,
		Rooms:   rooms,
	}
}

func timetableItemsToAPI(ds []domain.TimetableItem) []api.TimetableItem {
	result := make([]api.TimetableItem, len(ds))
	for i, d := range ds {
		result[i] = timetableItemToAPI(d)
	}
	return result
}

func toDomainTimetableItemFromRequest(req api.TimetableItemRequest) domain.TimetableItem {
	item := domain.TimetableItem{
		Subject: domain.Subject{ID: req.SubjectId},
	}
	if req.Slot != nil {
		dow := strings.TrimSpace(string(req.Slot.DayOfWeek))
		per := strings.TrimSpace(string(req.Slot.Period))
		if dow != "" && per != "" {
			item.Slot = &domain.TimetableSlot{
				DayOfWeek: domain.DayOfWeek(dow),
				Period:    domain.Period(per),
			}
		}
	}
	rooms := make([]domain.Room, len(req.RoomIds))
	for i, id := range req.RoomIds {
		rooms[i] = domain.Room{ID: id}
	}
	item.Rooms = rooms
	return item
}

// CourseRegistration

func courseRegistrationToAPI(d domain.CourseRegistration) api.CourseRegistration {
	return api.CourseRegistration{
		Id:      d.ID,
		UserId:  d.UserID,
		Subject: subjectToSummaryAPI(d.Subject),
	}
}

func courseRegistrationsToAPI(ds []domain.CourseRegistration) []api.CourseRegistration {
	result := make([]api.CourseRegistration, len(ds))
	for i, d := range ds {
		result[i] = courseRegistrationToAPI(d)
	}
	return result
}

func buildCourseRegistrationListFilter(params api.CourseRegistrationsV1ListParams) domain.CourseRegistrationListFilter {
	filter := domain.CourseRegistrationListFilter{
		UserID: params.UserId,
	}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	semesters := make([]domain.CourseSemester, len(params.Semesters))
	for i, s := range params.Semesters {
		semesters[i] = domain.CourseSemester(s)
	}
	filter.Semesters = semesters
	return filter
}

func toDomainCourseRegistrationFromRequest(req api.CourseRegistrationRequest) domain.CourseRegistration {
	return domain.CourseRegistration{
		UserID:  req.UserId,
		Subject: domain.Subject{ID: req.SubjectId},
	}
}
