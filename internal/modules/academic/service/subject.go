package service

import (
	"cmp"
	"context"
	"math"
	"slices"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
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
	subjects, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	if filter.SortByUserAttribute {
		sortSubjects(subjects, filter.SortCourse)
	}
	return subjects, nil
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

var courseOrder = map[domain.CourseType]int{
	domain.CourseTypeInformationSystem: 0,
	domain.CourseTypeInformationDesign: 1,
	domain.CourseTypeAdvancedICT:       2,
	domain.CourseTypeComplexSystem:     3,
	domain.CourseTypeIntelligentSystem: 4,
}

var gradeOrder = map[domain.Grade]int{
	domain.GradeB1: 0,
	domain.GradeB2: 1,
	domain.GradeB3: 2,
	domain.GradeB4: 3,
	domain.GradeM1: 4,
	domain.GradeM2: 5,
	domain.GradeD1: 6,
	domain.GradeD2: 7,
	domain.GradeD3: 8,
}

// courseRank はソートキーを返す。
// ユーザーのコースと一致する場合は -1（最優先）、
// それ以外は enum 定義順、Requirements が空なら math.MaxInt（最後）。
func courseRank(s domain.Subject, userCourse *domain.CourseType) int {
	if len(s.Requirements) == 0 {
		return math.MaxInt
	}
	best := math.MaxInt
	for _, r := range s.Requirements {
		if userCourse != nil && r.Course == *userCourse {
			return -1
		}
		if rank, ok := courseOrder[r.Course]; ok && rank < best {
			best = rank
		}
	}
	return best
}

// gradeRank はソートキーを返す。
// EligibleAttributes の中で最も若い学年の順位を返す。
// EligibleAttributes が空なら math.MaxInt（最後）。
func gradeRank(s domain.Subject) int {
	if len(s.EligibleAttributes) == 0 {
		return math.MaxInt
	}
	best := math.MaxInt
	for _, a := range s.EligibleAttributes {
		if rank, ok := gradeOrder[a.Grade]; ok && rank < best {
			best = rank
		}
	}
	return best
}

func sortSubjects(subjects []domain.Subject, userCourse *domain.CourseType) {
	slices.SortStableFunc(subjects, func(a, b domain.Subject) int {
		if c := cmp.Compare(courseRank(a, userCourse), courseRank(b, userCourse)); c != 0 {
			return c
		}
		return cmp.Compare(gradeRank(a), gradeRank(b))
	})
}
