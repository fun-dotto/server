package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type Subject struct {
	ID                      string                     `gorm:"type:uuid;primaryKey"`
	Name                    string                     `gorm:"not null"`
	Year                    int                        `gorm:"not null"`
	Semester                string                     `gorm:"not null"`
	Credit                  int                        `gorm:"not null"`
	Classification          string                     `gorm:"not null"`
	CulturalSubjectCategory string                     `gorm:"not null"`
	SyllabusID              string                     `gorm:"not null;uniqueIndex"`
	Syllabus                *Syllabus                  `gorm:"foreignKey:SyllabusID"`
	Faculties               []SubjectFaculty           `gorm:"foreignKey:SubjectID"`
	EligibleAttributes      []SubjectEligibleAttribute `gorm:"foreignKey:SubjectID"`
	Requirements            []SubjectRequirement       `gorm:"foreignKey:SubjectID"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type SubjectFaculty struct {
	ID        string   `gorm:"type:uuid;primaryKey"`
	SubjectID string   `gorm:"type:uuid;not null;index"`
	FacultyID string   `gorm:"type:uuid;not null"`
	Faculty   *Faculty `gorm:"foreignKey:FacultyID"`
	IsPrimary bool     `gorm:"not null"`
}

type SubjectEligibleAttribute struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	SubjectID string `gorm:"type:uuid;not null;index"`
	Grade     string `gorm:"not null"`
	Class     *string
}

type SubjectRequirement struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	SubjectID       string `gorm:"type:uuid;not null;index"`
	Course          string `gorm:"not null"`
	RequirementType string `gorm:"not null"`
}

func SubjectToDomain(m Subject) domain.Subject {
	faculties := make([]domain.SubjectFaculty, len(m.Faculties))
	for i, f := range m.Faculties {
		faculties[i] = domain.SubjectFaculty{
			Faculty:   FacultyToDomain(*f.Faculty),
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
		ID:                      m.ID,
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

func SubjectFromDomain(d domain.Subject) Subject {
	faculties := make([]SubjectFaculty, len(d.Faculties))
	for i, f := range d.Faculties {
		faculties[i] = SubjectFaculty{
			FacultyID: f.Faculty.ID,
			IsPrimary: f.IsPrimary,
		}
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
			Course:          string(r.Course),
			RequirementType: string(r.RequirementType),
		}
	}

	return Subject{
		ID:                      d.ID,
		Name:                    d.Name,
		Year:                    d.Year,
		Semester:                string(d.Semester),
		Credit:                  d.Credit,
		Classification:          string(d.Classification),
		CulturalSubjectCategory: string(d.CulturalSubjectCategory),
		SyllabusID:              d.SyllabusID,
		Faculties:               faculties,
		EligibleAttributes:      eligible,
		Requirements:            requirements,
	}
}
