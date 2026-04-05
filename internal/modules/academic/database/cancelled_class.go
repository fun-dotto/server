package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type CancelledClass struct {
	ID        string   `gorm:"type:uuid;primaryKey"`
	SubjectID string   `gorm:"type:uuid;not null;index"`
	Subject   *Subject `gorm:"foreignKey:SubjectID"`
	Date      string   `gorm:"type:date;not null;index"`
	Period    string   `gorm:"not null"`
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CancelledClassToDomain(m CancelledClass) domain.CancelledClass {
	date := m.Date
	if t, err := time.Parse("2006-01-02T15:04:05Z", m.Date); err == nil {
		date = t.Format("2006-01-02")
	} else if t, err := time.Parse(time.RFC3339, m.Date); err == nil {
		date = t.Format("2006-01-02")
	}
	d := domain.CancelledClass{
		ID:     m.ID,
		Date:   date,
		Period: domain.Period(m.Period),
	}
	if m.Comment != nil {
		d.Comment = *m.Comment
	}
	if m.Subject != nil {
		d.Subject = SubjectToDomain(*m.Subject)
	}
	return d
}

func CancelledClassFromDomain(d domain.CancelledClass) CancelledClass {
	var comment *string
	if d.Comment != "" {
		comment = &d.Comment
	}
	return CancelledClass{
		ID:        d.ID,
		SubjectID: d.Subject.ID,
		Date:      d.Date,
		Period:    string(d.Period),
		Comment:   comment,
	}
}
