package database

import (
	"time"

	"github.com/fun-dotto/subject-api/internal/domain"
)

type Subject struct {
	ID                      string                     `gorm:"type:uuid;primaryKey"`
	Name                    string                     `gorm:"type:text;not null"`
	FacultyID               string                     `gorm:"type:uuid;not null"`
	Faculty                 Faculty                    `gorm:"foreignKey:FacultyID"`
	Semester                string                     `gorm:"type:text;not null"`
	SyllabusID              string                     `gorm:"type:text;not null"`
	DayOfWeekTimetableSlots []DayOfWeekTimetableSlot   `gorm:"many2many:subject_day_of_week_timetable_slots;"`
	Categories              []SubjectCategory          `gorm:"many2many:subject_categories_subjects;"`
	EligibleAttributes      []SubjectEligibleAttribute `gorm:"foreignKey:SubjectID"`
	Requirements            []SubjectRequirement       `gorm:"foreignKey:SubjectID"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type SubjectEligibleAttribute struct {
	ID        string  `gorm:"type:uuid;primaryKey"`
	SubjectID string  `gorm:"type:uuid;not null;index"`
	Grade     string  `gorm:"type:text;not null"`
	Class     *string `gorm:"type:text"`
}

type SubjectRequirement struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	SubjectID       string `gorm:"type:uuid;not null;index"`
	CourseID        string `gorm:"type:uuid;not null"`
	Course          Course `gorm:"foreignKey:CourseID"`
	RequirementType string `gorm:"type:text;not null"`
}

func SubjectToDomain(m Subject) domain.Subject {
	slots := make([]domain.DayOfWeekTimetableSlot, len(m.DayOfWeekTimetableSlots))
	for i, s := range m.DayOfWeekTimetableSlots {
		slots[i] = DayOfWeekTimetableSlotToDomain(s)
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
			Course:          CourseToDomain(r.Course),
			RequirementType: domain.SubjectRequirementType(r.RequirementType),
		}
	}

	categories := make([]domain.SubjectCategory, len(m.Categories))
	for i, c := range m.Categories {
		categories[i] = SubjectCategoryToDomain(c)
	}

	return domain.Subject{
		ID:                      m.ID,
		Name:                    m.Name,
		Faculty:                 FacultyToDomain(m.Faculty),
		Semester:                domain.CourseSemester(m.Semester),
		DayOfWeekTimetableSlots: slots,
		EligibleAttributes:      eligible,
		Requirements:            requirements,
		Categories:              categories,
		SyllabusID:              m.SyllabusID,
	}
}

func SubjectFromDomain(d domain.Subject) Subject {
	slots := make([]DayOfWeekTimetableSlot, len(d.DayOfWeekTimetableSlots))
	for i, s := range d.DayOfWeekTimetableSlots {
		slots[i] = DayOfWeekTimetableSlotFromDomain(s)
	}

	eligible := make([]SubjectEligibleAttribute, len(d.EligibleAttributes))
	for i, e := range d.EligibleAttributes {
		attr := SubjectEligibleAttribute{
			Grade: string(e.Grade),
		}
		if e.Class != nil {
			c := string(*e.Class)
			attr.Class = &c
		}
		eligible[i] = attr
	}

	requirements := make([]SubjectRequirement, len(d.Requirements))
	for i, r := range d.Requirements {
		requirements[i] = SubjectRequirement{
			CourseID:        r.Course.ID,
			Course:          CourseFromDomain(r.Course),
			RequirementType: string(r.RequirementType),
		}
	}

	categories := make([]SubjectCategory, len(d.Categories))
	for i, c := range d.Categories {
		categories[i] = SubjectCategoryFromDomain(c)
	}

	return Subject{
		ID:                      d.ID,
		Name:                    d.Name,
		FacultyID:               d.Faculty.ID,
		Faculty:                 FacultyFromDomain(d.Faculty),
		Semester:                string(d.Semester),
		SyllabusID:              d.SyllabusID,
		DayOfWeekTimetableSlots: slots,
		Categories:              categories,
		EligibleAttributes:      eligible,
		Requirements:            requirements,
	}
}
