package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type subjectRepository interface {
	List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	GetBySyllabusID(ctx context.Context, syllabusID string) (domain.Subject, error)
	Upsert(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
}

type syllabusRepository interface {
	GetByID(ctx context.Context, id string) (domain.Syllabus, error)
}

type SubjectService struct {
	repo         subjectRepository
	syllabusRepo syllabusRepository
}

func NewSubjectService(repo subjectRepository, syllabusRepo syllabusRepository) *SubjectService {
	return &SubjectService{repo: repo, syllabusRepo: syllabusRepo}
}

func (s *SubjectService) List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error) {
	return s.repo.List(ctx, filter)
}

func (s *SubjectService) GetByID(ctx context.Context, id string) (domain.Subject, error) {
	return s.repo.GetByID(ctx, id)
}

// Upsert はシラバスIDからシラバスを取得し、Subjectを導出して作成または更新する。
// TODO: Syllabus -> Subject 導出ロジックの詳細実装。現在はスタブ。
func (s *SubjectService) Upsert(ctx context.Context, syllabusID string) (domain.Subject, error) {
	syllabus, err := s.syllabusRepo.GetByID(ctx, syllabusID)
	if err != nil {
		return domain.Subject{}, err
	}

	subject := deriveSubjectFromSyllabus(syllabus)
	return s.repo.Upsert(ctx, subject)
}

func (s *SubjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubjectService) GetSyllabus(ctx context.Context, subjectID string) (domain.Syllabus, error) {
	subject, err := s.repo.GetByID(ctx, subjectID)
	if err != nil {
		return domain.Syllabus{}, err
	}
	return s.syllabusRepo.GetByID(ctx, subject.SyllabusID)
}

// deriveSubjectFromSyllabus はシラバスからSubjectドメインモデルを導出する。
// TODO: Year, Semester が未設定（ゼロ値）のまま。Syllabus にこれらの情報がないため、導出元の追加か Upsert の引数追加が必要。
// TODO: grades, targetCourses, classifications等のテキストフィールドからEligibleAttributes, Requirements等を構造化する。
func deriveSubjectFromSyllabus(syllabus domain.Syllabus) domain.Subject {
	return domain.Subject{
		Name:               syllabus.Name,
		Credit:             syllabus.Credit,
		SyllabusID:         syllabus.ID,
		Faculties:          []domain.SubjectFaculty{},
		EligibleAttributes: []domain.SubjectTargetClass{},
		Requirements:       []domain.SubjectRequirement{},
	}
}
