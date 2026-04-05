package database

import (
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type CourseRegistration struct {
	ID        string   `gorm:"type:uuid;primaryKey"`
	// TODO: (user_id, subject_id) の複合ユニーク制約を追加する（既存データの重複確認・クリーンアップ後に uniqueIndex:idx_user_subject へ変更）
	UserID    string   `gorm:"not null;index"`
	SubjectID string   `gorm:"type:uuid;not null;index"`
	Subject   *Subject `gorm:"foreignKey:SubjectID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CourseRegistrationToDomain(m CourseRegistration) domain.CourseRegistration {
	var subject domain.Subject
	if m.Subject != nil {
		subject = SubjectToDomain(*m.Subject)
	}

	return domain.CourseRegistration{
		ID:      m.ID,
		UserID:  m.UserID,
		Subject: subject,
	}
}

func CourseRegistrationFromDomain(d domain.CourseRegistration) CourseRegistration {
	return CourseRegistration{
		ID:        d.ID,
		UserID:    d.UserID,
		SubjectID: d.Subject.ID,
	}
}
