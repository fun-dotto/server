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
		sortSubjects(subjects, filter.SortCourse, filter.SortGrade)
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

var classificationOrder = map[domain.SubjectClassification]int{
	domain.SubjectClassificationSpecialized:         0,
	domain.SubjectClassificationResearchInstruction: 1,
	domain.SubjectClassificationCultural:            2,
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

// classificationRank は科目区分のソートキーを返す。
// 専門 → 研究指導 → 教養 の順。
func classificationRank(s domain.Subject) int {
	if rank, ok := classificationOrder[s.Classification]; ok {
		return rank
	}
	return math.MaxInt
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
// ユーザーの学年と一致する場合は -1（最優先）、
// それ以外は学年の昇順（B1=0, B2=1, ...）、EligibleAttributes が空なら math.MaxInt（最後）。
func gradeRank(s domain.Subject, userGrade *domain.Grade) int {
	if len(s.EligibleAttributes) == 0 {
		return math.MaxInt
	}
	best := math.MaxInt
	for _, a := range s.EligibleAttributes {
		if userGrade != nil && a.Grade == *userGrade {
			return -1
		}
		if rank, ok := gradeOrder[a.Grade]; ok && rank < best {
			best = rank
		}
	}
	return best
}

func sortSubjects(subjects []domain.Subject, userCourse *domain.CourseType, userGrade *domain.Grade) {
	slices.SortStableFunc(subjects, func(a, b domain.Subject) int {
		if c := cmp.Compare(classificationRank(a), classificationRank(b)); c != 0 {
			return c
		}
		if c := cmp.Compare(courseRank(a, userCourse), courseRank(b, userCourse)); c != 0 {
			return c
		}
		return cmp.Compare(gradeRank(a, userGrade), gradeRank(b, userGrade))
	})
}
