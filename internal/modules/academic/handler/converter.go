package handler

import (
	api "github.com/fun-dotto/subject-api/generated"
	"github.com/fun-dotto/subject-api/internal/domain"
)

// Course

func courseToAPI(d domain.Course) api.Course {
	return api.Course{
		Id:   d.ID,
		Name: d.Name,
	}
}

func coursesToAPI(ds []domain.Course) []api.Course {
	result := make([]api.Course, len(ds))
	for i, d := range ds {
		result[i] = courseToAPI(d)
	}
	return result
}

func courseRequestToDomain(r api.CourseRequest) domain.Course {
	return domain.Course{
		Name: r.Name,
	}
}

// Faculty

func facultyToAPI(d domain.Faculty) api.Faculty {
	return api.Faculty{
		Id:    d.ID,
		Name:  d.Name,
		Email: d.Email,
	}
}

func facultiesToAPI(ds []domain.Faculty) []api.Faculty {
	result := make([]api.Faculty, len(ds))
	for i, d := range ds {
		result[i] = facultyToAPI(d)
	}
	return result
}

func facultyRequestToDomain(r api.FacultyRequest) domain.Faculty {
	return domain.Faculty{
		Name:  r.Name,
		Email: r.Email,
	}
}

// DayOfWeekTimetableSlot

func slotToAPI(d domain.DayOfWeekTimetableSlot) api.DayOfWeekTimetableSlot {
	return api.DayOfWeekTimetableSlot{
		Id:            d.ID,
		DayOfWeek:     api.DottoFoundationV1DayOfWeek(d.DayOfWeek),
		TimetableSlot: api.DottoFoundationV1TimetableSlot(d.TimetableSlot),
	}
}

func slotsToAPI(ds []domain.DayOfWeekTimetableSlot) []api.DayOfWeekTimetableSlot {
	result := make([]api.DayOfWeekTimetableSlot, len(ds))
	for i, d := range ds {
		result[i] = slotToAPI(d)
	}
	return result
}

func slotRequestToDomain(r api.DayOfWeekTimetableSlotRequest) domain.DayOfWeekTimetableSlot {
	return domain.DayOfWeekTimetableSlot{
		DayOfWeek:     domain.DayOfWeek(r.DayOfWeek),
		TimetableSlot: domain.TimetableSlot(r.TimetableSlot),
	}
}

// SubjectCategory

func subjectCategoryToAPI(d domain.SubjectCategory) api.SubjectCategory {
	return api.SubjectCategory{
		Id:   d.ID,
		Name: d.Name,
	}
}

func subjectCategoriesToAPI(ds []domain.SubjectCategory) []api.SubjectCategory {
	result := make([]api.SubjectCategory, len(ds))
	for i, d := range ds {
		result[i] = subjectCategoryToAPI(d)
	}
	return result
}

func subjectCategoryRequestToDomain(r api.SubjectCategoryRequest) domain.SubjectCategory {
	return domain.SubjectCategory{
		Name: r.Name,
	}
}

// Subject

func subjectToAPI(d domain.Subject) api.Subject {
	slots := make([]api.DayOfWeekTimetableSlot, len(d.DayOfWeekTimetableSlots))
	for i, s := range d.DayOfWeekTimetableSlots {
		slots[i] = slotToAPI(s)
	}

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
			Course:          courseToAPI(r.Course),
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}

	categories := make([]api.SubjectCategory, len(d.Categories))
	for i, c := range d.Categories {
		categories[i] = subjectCategoryToAPI(c)
	}

	return api.Subject{
		Id:                      d.ID,
		Name:                    d.Name,
		Faculty:                 facultyToAPI(d.Faculty),
		Semester:                api.DottoFoundationV1CourseSemester(d.Semester),
		DayOfWeekTimetableSlots: slots,
		EligibleAttributes:      eligible,
		Requirements:            requirements,
		Categories:              categories,
		SyllabusId:              d.SyllabusID,
	}
}

func subjectsToAPI(ds []domain.Subject) []api.Subject {
	result := make([]api.Subject, len(ds))
	for i, d := range ds {
		result[i] = subjectToAPI(d)
	}
	return result
}

func subjectRequestToDomain(r api.SubjectRequest) domain.Subject {
	slots := make([]domain.DayOfWeekTimetableSlot, len(r.DayOfWeekTimetableSlotIds))
	for i, id := range r.DayOfWeekTimetableSlotIds {
		slots[i] = domain.DayOfWeekTimetableSlot{ID: id}
	}

	eligible := make([]domain.SubjectTargetClass, len(r.EligibleAttributes))
	for i, e := range r.EligibleAttributes {
		tc := domain.SubjectTargetClass{
			Grade: domain.Grade(e.Grade),
		}
		if e.Class != nil {
			c := domain.Class(*e.Class)
			tc.Class = &c
		}
		eligible[i] = tc
	}

	requirements := make([]domain.SubjectRequirement, len(r.Requirements))
	for i, req := range r.Requirements {
		requirements[i] = domain.SubjectRequirement{
			Course:          domain.Course{ID: req.CourseId},
			RequirementType: domain.SubjectRequirementType(req.RequirementType),
		}
	}

	categories := make([]domain.SubjectCategory, len(r.CategoryIds))
	for i, id := range r.CategoryIds {
		categories[i] = domain.SubjectCategory{ID: id}
	}

	return domain.Subject{
		Name:                    r.Name,
		Faculty:                 domain.Faculty{ID: r.FacultyId},
		Semester:                domain.CourseSemester(r.Semester),
		DayOfWeekTimetableSlots: slots,
		EligibleAttributes:      eligible,
		Requirements:            requirements,
		Categories:              categories,
		SyllabusID:              r.SyllabusId,
	}
}
