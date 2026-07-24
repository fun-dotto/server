package service

import (
	"math"
	"testing"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func coursePtr(c domain.CourseType) *domain.CourseType { return &c }

func gradePtr(g domain.Grade) *domain.Grade { return &g }

func TestClassificationRank(t *testing.T) {
	tests := []struct {
		name    string
		subject domain.Subject
		want    int
	}{
		{
			name:    "専門が最優先",
			subject: domain.Subject{Classification: domain.SubjectClassificationSpecialized},
			want:    0,
		},
		{
			name:    "研究指導は専門の次",
			subject: domain.Subject{Classification: domain.SubjectClassificationResearchInstruction},
			want:    1,
		},
		{
			name:    "教養は最後",
			subject: domain.Subject{Classification: domain.SubjectClassificationCultural},
			want:    2,
		},
		{
			name:    "未知の区分は最後尾",
			subject: domain.Subject{Classification: domain.SubjectClassification("Unknown")},
			want:    math.MaxInt,
		},
		{
			name:    "空文字も最後尾",
			subject: domain.Subject{},
			want:    math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classificationRank(tt.subject); got != tt.want {
				t.Errorf("classificationRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourseRank(t *testing.T) {
	subjectWith := func(courses ...domain.CourseType) domain.Subject {
		reqs := make([]domain.SubjectRequirement, len(courses))
		for i, c := range courses {
			reqs[i] = domain.SubjectRequirement{Course: c}
		}
		return domain.Subject{Requirements: reqs}
	}

	tests := []struct {
		name       string
		subject    domain.Subject
		userCourse *domain.CourseType
		want       int
	}{
		{
			name:       "ユーザーのコースと一致すれば最優先",
			subject:    subjectWith(domain.CourseTypeComplexSystem),
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			want:       -1,
		},
		{
			name:       "複数コースのうち1つでも一致すれば最優先",
			subject:    subjectWith(domain.CourseTypeInformationSystem, domain.CourseTypeIntelligentSystem),
			userCourse: coursePtr(domain.CourseTypeIntelligentSystem),
			want:       -1,
		},
		{
			name:       "一致しなければ enum 定義順のうち最小",
			subject:    subjectWith(domain.CourseTypeIntelligentSystem, domain.CourseTypeInformationDesign),
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			want:       1, // InformationDesign
		},
		{
			name:       "ユーザーのコースが未設定なら enum 定義順のみ",
			subject:    subjectWith(domain.CourseTypeAdvancedICT),
			userCourse: nil,
			want:       2,
		},
		{
			name:       "Requirements が空なら最後尾",
			subject:    domain.Subject{},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			want:       math.MaxInt,
		},
		{
			name:       "未知のコースのみなら最後尾",
			subject:    subjectWith(domain.CourseType("Unknown")),
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			want:       math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := courseRank(tt.subject, tt.userCourse); got != tt.want {
				t.Errorf("courseRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGradeRank(t *testing.T) {
	subjectWith := func(grades ...domain.Grade) domain.Subject {
		attrs := make([]domain.SubjectTargetClass, len(grades))
		for i, g := range grades {
			attrs[i] = domain.SubjectTargetClass{Grade: g}
		}
		return domain.Subject{EligibleAttributes: attrs}
	}

	tests := []struct {
		name      string
		subject   domain.Subject
		userGrade *domain.Grade
		want      int
	}{
		{
			name:      "ユーザーの学年と一致すれば最優先",
			subject:   subjectWith(domain.GradeB3),
			userGrade: gradePtr(domain.GradeB3),
			want:      -1,
		},
		{
			name:      "複数学年のうち1つでも一致すれば最優先",
			subject:   subjectWith(domain.GradeB1, domain.GradeB4),
			userGrade: gradePtr(domain.GradeB4),
			want:      -1,
		},
		{
			name:      "一致しなければ最も低い学年",
			subject:   subjectWith(domain.GradeM1, domain.GradeB2),
			userGrade: gradePtr(domain.GradeB4),
			want:      1, // B2
		},
		{
			name:      "ユーザーの学年が未設定なら学年昇順のみ",
			subject:   subjectWith(domain.GradeD1),
			userGrade: nil,
			want:      6,
		},
		{
			name:      "EligibleAttributes が空なら最後尾",
			subject:   domain.Subject{},
			userGrade: gradePtr(domain.GradeB1),
			want:      math.MaxInt,
		},
		{
			name:      "未知の学年のみなら最後尾",
			subject:   subjectWith(domain.Grade("B9")),
			userGrade: gradePtr(domain.GradeB1),
			want:      math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gradeRank(tt.subject, tt.userGrade); got != tt.want {
				t.Errorf("gradeRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

// newSubject はソート対象の科目を組み立てるヘルパー。
func newSubject(id string, classification domain.SubjectClassification, course domain.CourseType, grade domain.Grade) domain.Subject {
	return domain.Subject{
		ID:                 id,
		Classification:     classification,
		Requirements:       []domain.SubjectRequirement{{Course: course}},
		EligibleAttributes: []domain.SubjectTargetClass{{Grade: grade}},
	}
}

func TestSortSubjects(t *testing.T) {
	tests := []struct {
		name       string
		subjects   []domain.Subject
		userCourse *domain.CourseType
		userGrade  *domain.Grade
		wantIDs    []string
	}{
		{
			name: "第1キー: 専門 → 研究指導 → 教養",
			subjects: []domain.Subject{
				newSubject("cultural", domain.SubjectClassificationCultural, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("research", domain.SubjectClassificationResearchInstruction, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("specialized", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB1),
			wantIDs:    []string{"specialized", "research", "cultural"},
		},
		{
			name: "第2キー: 同一区分ならユーザーのコースが最優先",
			subjects: []domain.Subject{
				newSubject("other", domain.SubjectClassificationSpecialized, domain.CourseTypeIntelligentSystem, domain.GradeB1),
				newSubject("mine", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB1),
			wantIDs:    []string{"mine", "other"},
		},
		{
			name: "第3キー: 区分もコースも同じならユーザーの学年が最優先",
			subjects: []domain.Subject{
				newSubject("b1", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("b3", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB3),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB3),
			wantIDs:    []string{"b3", "b1"},
		},
		{
			name: "区分がコースより優先される（教養かつ一致コースは専門の後）",
			subjects: []domain.Subject{
				newSubject("cultural-mine", domain.SubjectClassificationCultural, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("specialized-other", domain.SubjectClassificationSpecialized, domain.CourseTypeIntelligentSystem, domain.GradeB4),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB1),
			wantIDs:    []string{"specialized-other", "cultural-mine"},
		},
		{
			name: "Requirements / EligibleAttributes が空の科目は最後尾",
			subjects: []domain.Subject{
				{ID: "empty", Classification: domain.SubjectClassificationSpecialized},
				newSubject("filled", domain.SubjectClassificationSpecialized, domain.CourseTypeIntelligentSystem, domain.GradeB1),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB4),
			wantIDs:    []string{"filled", "empty"},
		},
		{
			name: "全キー同順位なら元の順序を維持する（stable）",
			subjects: []domain.Subject{
				newSubject("first", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("second", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
				newSubject("third", domain.SubjectClassificationSpecialized, domain.CourseTypeComplexSystem, domain.GradeB1),
			},
			userCourse: coursePtr(domain.CourseTypeComplexSystem),
			userGrade:  gradePtr(domain.GradeB1),
			wantIDs:    []string{"first", "second", "third"},
		},
		{
			name: "コース・学年が未設定のユーザーでも区分と定義順で並ぶ",
			subjects: []domain.Subject{
				newSubject("cultural", domain.SubjectClassificationCultural, domain.CourseTypeInformationSystem, domain.GradeB1),
				newSubject("specialized", domain.SubjectClassificationSpecialized, domain.CourseTypeIntelligentSystem, domain.GradeB4),
			},
			userCourse: nil,
			userGrade:  nil,
			wantIDs:    []string{"specialized", "cultural"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortSubjects(tt.subjects, tt.userCourse, tt.userGrade)

			got := make([]string, len(tt.subjects))
			for i, s := range tt.subjects {
				got[i] = s.ID
			}
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("sortSubjects() = %v, want %v", got, tt.wantIDs)
			}
			for i := range got {
				if got[i] != tt.wantIDs[i] {
					t.Errorf("sortSubjects() = %v, want %v", got, tt.wantIDs)
					break
				}
			}
		})
	}
}
