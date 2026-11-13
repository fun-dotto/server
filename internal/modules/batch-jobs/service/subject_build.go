package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/fun-dotto/server/internal/shared/subjectbuild"
	"gorm.io/gorm"
)

// SubjectBuildSummary は科目生成ジョブの結果。
type SubjectBuildSummary struct {
	Created int
	Skipped int
	Failed  int
}

type SubjectBuildService struct {
	db *gorm.DB
}

func NewSubjectBuildService(db *gorm.DB) *SubjectBuildService {
	return &SubjectBuildService{db: db}
}

// Build は科目が存在しない year の授業について、最新バージョンのシラバスから科目と子テーブルを作成する。
// 既存の科目は更新しない（紐付いているバージョンを使い続ける）。
func (s *SubjectBuildService) Build(ctx context.Context, year int) (SubjectBuildSummary, error) {
	var summary SubjectBuildSummary
	db := s.db.WithContext(ctx)

	var syllabuses []model.Syllabus
	err := db.Raw(`
SELECT DISTINCT ON (sy.lesson_id) sy.*
FROM syllabuses sy
WHERE sy.year = ?
  AND NOT EXISTS (
    SELECT 1 FROM subjects s
    JOIN syllabuses linked ON linked.id = s.syllabus_id
    WHERE linked.lesson_id = sy.lesson_id AND linked.year = sy.year
  )
ORDER BY sy.lesson_id, sy.fetched_at DESC, sy.created_at DESC`, year).Scan(&syllabuses).Error
	if err != nil {
		return summary, fmt.Errorf("find syllabuses without subject: %w", err)
	}

	var faculties []model.Faculty
	if err := db.Find(&faculties).Error; err != nil {
		return summary, fmt.Errorf("list faculties: %w", err)
	}

	for _, sy := range syllabuses {
		res, err := subjectbuild.Build(sy)
		if err != nil {
			summary.Skipped++
			slog.WarnContext(ctx, "skip subject", "lesson_id", sy.LessonID, "year", sy.Year, "reason", err.Error())
			continue
		}
		for _, w := range res.Warnings {
			slog.WarnContext(ctx, "subject build warning", "lesson_id", sy.LessonID, "year", sy.Year, "warning", w)
		}
		if missing := res.ResolveFaculties(faculties); len(missing) > 0 {
			slog.WarnContext(ctx, "faculties not found", "lesson_id", sy.LessonID, "year", sy.Year, "names", missing)
		}
		subject := res.Subject
		if err := db.Omit("Syllabus", "Faculties.Faculty").Create(&subject).Error; err != nil {
			summary.Failed++
			slog.ErrorContext(ctx, "failed to create subject", "lesson_id", sy.LessonID, "year", sy.Year, "error", err)
			continue
		}
		summary.Created++
	}
	if summary.Failed > 0 {
		return summary, fmt.Errorf("%d subjects failed", summary.Failed)
	}
	return summary, nil
}
