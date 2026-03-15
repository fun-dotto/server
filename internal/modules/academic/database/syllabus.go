package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type Syllabus struct {
	ID                           string `gorm:"primaryKey"`
	Name                         string `gorm:"not null"`
	EnName                       string `gorm:"not null"`
	Grades                       string `gorm:"not null"`
	Credit                       int    `gorm:"not null"`
	FacultyNames                 string `gorm:"not null"`
	PracticalHomeFacultyCategory string `gorm:"not null"`
	MultiplePersonTeachingForm   string `gorm:"not null"`
	TeachingForm                 string `gorm:"not null"`
	Summary                      string `gorm:"not null"`
	LearningOutcomes             string `gorm:"not null"`
	Assignments                  string `gorm:"not null"`
	EvaluationMethod             string `gorm:"not null"`
	Textbooks                    string `gorm:"not null"`
	ReferenceBooks               string `gorm:"not null"`
	Prerequisites                string `gorm:"not null"`
	PreLearning                  string `gorm:"not null"`
	PostLearning                 string `gorm:"not null"`
	Notes                        string `gorm:"not null"`
	Keywords                     string `gorm:"not null"`
	TargetCourses                string `gorm:"not null"`
	TargetAreas                  string `gorm:"not null"`
	Classifications              string `gorm:"not null"`
	TeachingLanguage             string `gorm:"not null"`
	ContentsAndSchedule          string `gorm:"not null"`
	TeachingAndExamForm          string `gorm:"not null"`
	DspoSubject                  string `gorm:"not null"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

func SyllabusToDomain(m Syllabus) domain.Syllabus {
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
		DspoSubject:                  m.DspoSubject,
	}
}

func SyllabusFromDomain(d domain.Syllabus) Syllabus {
	return Syllabus{
		ID:                           d.ID,
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
		DspoSubject:                  d.DspoSubject,
	}
}
