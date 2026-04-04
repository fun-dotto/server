package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type subjectRepository interface {
	List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
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

