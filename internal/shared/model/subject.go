package model

import "github.com/google/uuid"

type Subject struct {
	Common

	Name                    string                     `gorm:"not null"`
	Year                    int                        `gorm:"not null"`
	Semester                string                     `gorm:"not null"`
	Credit                  int                        `gorm:"not null"`
	Classification          string                     `gorm:"not null"`
	CulturalSubjectCategory string                     `gorm:"not null"`
	SyllabusID              string                     `gorm:"not null;uniqueIndex"`
	Syllabus                *Syllabus                  `gorm:"foreignKey:SyllabusID;constraint:OnUpdate:CASCADE"`
	Faculties               []SubjectFaculty           `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`
	EligibleAttributes      []SubjectEligibleAttribute `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`
	Requirements            []SubjectRequirement       `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE"`
}

type SubjectFaculty struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	FacultyID uuid.UUID `gorm:"type:uuid;not null"`
	Faculty   *Faculty  `gorm:"foreignKey:FacultyID"`
	IsPrimary bool      `gorm:"not null"`
}

type SubjectEligibleAttribute struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Grade     string    `gorm:"not null"`
	Class     *string
}

type SubjectRequirement struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SubjectID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Course          string    `gorm:"not null"`
	RequirementType string    `gorm:"not null"`
}
