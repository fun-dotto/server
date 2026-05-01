package model

import "time"

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
	DsopSubject                  string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
