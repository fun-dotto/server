package service

import (
	"context"
	"fmt"
	"log/slog"
	"uuid"

	"github.com/fun-dotto/server/internal/shared/model"
)

type SubjectBuildRepository interface {
	// ListSyllabusesWithoutSubject は年度のシラバスのうち、同じ (lesson_id, year) のどのバージョンにも科目が紐付いていないものを、(lesson_id, year) ごとに最新の 1 件ずつ返す。
	ListSyllabusesWithoutSubject(ctx context.Context, year int) ([]model.Syllabus, error)
	ListFaculties(ctx context.Context) ([]model.Faculty, error)
	// CreateSubject は科目と子テーブルを 1 トランザクションで作成する。
	CreateSubject(ctx context.Context, subject *model.Subject) error
}

type SubjectBuildService struct {
	repo SubjectBuildRepository
}

func NewSubjectBuildService(repo SubjectBuildRepository) *SubjectBuildService {
	return &SubjectBuildService{repo: repo}
}

type SubjectBuildSummary struct {
	Created int
	Skipped int
	Failed  int
}

// BuildSubjects は年度のシラバスのうち科目が存在しないものから科目を新規作成する。既存の科目は更新しない。
// 1 科目の作成に失敗しても他の科目の作成は続ける。
func (s *SubjectBuildService) BuildSubjects(ctx context.Context, year int) (SubjectBuildSummary, error) {
	var summary SubjectBuildSummary

	syllabuses, err := s.repo.ListSyllabusesWithoutSubject(ctx, year)
	if err != nil {
		return summary, fmt.Errorf("list syllabuses without subject: %w", err)
	}
	faculties, err := s.repo.ListFaculties(ctx)
	if err != nil {
		return summary, fmt.Errorf("list faculties: %w", err)
	}
	facultyIDs := make(map[string]uuid.UUID, len(faculties))
	for _, f := range faculties {
		name := NormalizeFacultyName(f.Name)
		if _, dup := facultyIDs[name]; dup {
			// 同名の教員が複数いる場合はどちらか決められないため、ゼロ値にして紐付けない
			slog.WarnContext(ctx, "duplicate faculty name", "faculty_name", name)
			facultyIDs[name] = uuid.UUID{}
			continue
		}
		facultyIDs[name] = f.ID
	}

	for _, syl := range syllabuses {
		structured, ok := StructureSyllabus(syl)
		if !ok {
			slog.WarnContext(ctx, "skip syllabus: unknown term", "lesson_id", syl.LessonID, "year", syl.Year, "term", syl.Term)
			summary.Skipped++
			continue
		}
		subject := buildSubjectModel(ctx, syl, structured, facultyIDs)
		if err := s.repo.CreateSubject(ctx, &subject); err != nil {
			slog.ErrorContext(ctx, "failed to create subject", "lesson_id", syl.LessonID, "year", syl.Year, "error", err)
			summary.Failed++
			continue
		}
		summary.Created++
	}
	return summary, nil
}

func buildSubjectModel(ctx context.Context, syl model.Syllabus, st StructuredSubject, facultyIDs map[string]uuid.UUID) model.Subject {
	var faculties []model.SubjectFaculty
	for _, f := range st.FacultyNames {
		id, ok := facultyIDs[f.Name]
		if ok && id == (uuid.UUID{}) {
			slog.WarnContext(ctx, "faculty name is ambiguous", "lesson_id", syl.LessonID, "year", syl.Year, "faculty_name", f.Name)
			continue
		}
		if !ok {
			slog.WarnContext(ctx, "faculty not found", "lesson_id", syl.LessonID, "year", syl.Year, "faculty_name", f.Name)
			continue
		}
		faculties = append(faculties, model.SubjectFaculty{FacultyID: id, IsPrimary: f.IsPrimary})
	}
	return model.Subject{
		Name:                    st.Name,
		Year:                    st.Year,
		Semester:                string(st.Semester),
		Credit:                  st.Credit,
		Classification:          string(st.Classification),
		CulturalSubjectCategory: st.CulturalSubjectCategory,
		SyllabusID:              syl.ID,
		Faculties:               faculties,
		EligibleAttributes:      st.EligibleAttributes,
		Requirements:            st.Requirements,
	}
}
