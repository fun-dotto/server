package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type SubjectBuildRepository struct {
	db *gorm.DB
}

func NewSubjectBuildRepository(db *gorm.DB) *SubjectBuildRepository {
	return &SubjectBuildRepository{db: db}
}

func (r *SubjectBuildRepository) ListSyllabusesWithoutSubject(ctx context.Context, year int) ([]model.Syllabus, error) {
	var syllabuses []model.Syllabus
	// 科目の有無は (lesson_id, year) 単位で判定する。旧バージョンに科目が紐付いていれば、新しいバージョンから重複して作らない。
	err := r.db.WithContext(ctx).Raw(`
SELECT DISTINCT ON (s.lesson_id) s.*
FROM syllabuses s
WHERE s.year = ?
  AND NOT EXISTS (
    SELECT 1
    FROM subjects sub
    JOIN syllabuses linked ON linked.id = sub.syllabus_id
    WHERE linked.lesson_id = s.lesson_id AND linked.year = s.year
  )
ORDER BY s.lesson_id ASC, s.fetched_at DESC, s.id ASC`, year).Scan(&syllabuses).Error
	return syllabuses, err
}

func (r *SubjectBuildRepository) ListFaculties(ctx context.Context) ([]model.Faculty, error) {
	var faculties []model.Faculty
	err := r.db.WithContext(ctx).Order("id ASC").Find(&faculties).Error
	return faculties, err
}

func (r *SubjectBuildRepository) CreateSubject(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(subject).Error
	})
}
