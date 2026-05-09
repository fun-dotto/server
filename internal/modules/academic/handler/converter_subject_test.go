package handler

import (
	"testing"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestSubjectToAPI(t *testing.T) {
	classA := domain.ClassA
	got := subjectToAPI(domain.Subject{
		ID:       "subj-1",
		Name:     "アルゴリズム論",
		Year:     2026,
		Semester: domain.CourseSemesterQ1,
		Credit:   2,
		Faculties: []domain.SubjectFaculty{
			{
				Faculty:   domain.Faculty{ID: "fac-1", Name: "教員A", Email: "a@example.com"},
				IsPrimary: true,
			},
		},
		EligibleAttributes: []domain.SubjectTargetClass{
			{Grade: domain.GradeB2, Class: &classA},
		},
		Requirements: []domain.SubjectRequirement{
			{Course: domain.CourseTypeInformationSystem, RequirementType: domain.SubjectRequirementTypeRequired},
		},
	})

	if got.Id != "subj-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "subj-1")
	}
	if got.Name != "アルゴリズム論" {
		t.Errorf("Name: got %q, want %q", got.Name, "アルゴリズム論")
	}
	if got.Year != 2026 {
		t.Errorf("Year: got %d, want %d", got.Year, 2026)
	}
	if got.Semester != api.Q1 {
		t.Errorf("Semester: got %q, want %q", got.Semester, api.Q1)
	}
	if got.Credit != 2 {
		t.Errorf("Credit: got %d, want %d", got.Credit, 2)
	}
	if len(got.Faculties) != 1 {
		t.Fatalf("Faculties len: got %d, want 1", len(got.Faculties))
	}
	if got.Faculties[0].Faculty.Id != "fac-1" {
		t.Errorf("Faculties[0].Faculty.Id: got %q, want %q", got.Faculties[0].Faculty.Id, "fac-1")
	}
	if !got.Faculties[0].IsPrimary {
		t.Error("Faculties[0].IsPrimary: got false, want true")
	}
	if got.EligibleAttributes != nil {
		t.Errorf("EligibleAttributes: got %+v, want nil (not included in subjectToAPI)", got.EligibleAttributes)
	}
	if got.Requirements != nil {
		t.Errorf("Requirements: got %+v, want nil (not included in subjectToAPI)", got.Requirements)
	}
}

func TestSubjectToDetailAPI(t *testing.T) {
	classA := domain.ClassA
	got := subjectToDetailAPI(domain.Subject{
		ID:       "subj-1",
		Name:     "アルゴリズム論",
		Year:     2026,
		Semester: domain.CourseSemesterQ1,
		Credit:   2,
		Faculties: []domain.SubjectFaculty{
			{
				Faculty:   domain.Faculty{ID: "fac-1", Name: "教員A", Email: "a@example.com"},
				IsPrimary: true,
			},
		},
		EligibleAttributes: []domain.SubjectTargetClass{
			{Grade: domain.GradeB2, Class: &classA},
			{Grade: domain.GradeM1, Class: nil},
		},
		Requirements: []domain.SubjectRequirement{
			{Course: domain.CourseTypeInformationSystem, RequirementType: domain.SubjectRequirementTypeRequired},
			{Course: domain.CourseTypeInformationDesign, RequirementType: domain.SubjectRequirementTypeOptional},
		},
	})

	if got.Id != "subj-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "subj-1")
	}
	if got.Name != "アルゴリズム論" {
		t.Errorf("Name: got %q, want %q", got.Name, "アルゴリズム論")
	}
	if got.Semester != api.Q1 {
		t.Errorf("Semester: got %q, want %q", got.Semester, api.Q1)
	}

	if got.EligibleAttributes == nil {
		t.Fatal("EligibleAttributes: got nil, want non-nil")
	}
	ea := *got.EligibleAttributes
	if len(ea) != 2 {
		t.Fatalf("EligibleAttributes len: got %d, want 2", len(ea))
	}
	if ea[0].Grade != api.B2 {
		t.Errorf("EligibleAttributes[0].Grade: got %q, want %q", ea[0].Grade, api.B2)
	}
	if ea[0].Class == nil {
		t.Fatal("EligibleAttributes[0].Class: got nil, want non-nil")
	}
	if *ea[0].Class != api.A {
		t.Errorf("EligibleAttributes[0].Class: got %q, want %q", *ea[0].Class, api.A)
	}
	if ea[1].Grade != api.M1 {
		t.Errorf("EligibleAttributes[1].Grade: got %q, want %q", ea[1].Grade, api.M1)
	}
	if ea[1].Class != nil {
		t.Errorf("EligibleAttributes[1].Class: got %+v, want nil", ea[1].Class)
	}

	if got.Requirements == nil {
		t.Fatal("Requirements: got nil, want non-nil")
	}
	reqs := *got.Requirements
	if len(reqs) != 2 {
		t.Fatalf("Requirements len: got %d, want 2", len(reqs))
	}
	if reqs[0].Course != api.InformationSystem {
		t.Errorf("Requirements[0].Course: got %q, want %q", reqs[0].Course, api.InformationSystem)
	}
	if reqs[0].RequirementType != api.Required {
		t.Errorf("Requirements[0].RequirementType: got %q, want %q", reqs[0].RequirementType, api.Required)
	}
	if reqs[1].Course != api.InformationDesign {
		t.Errorf("Requirements[1].Course: got %q, want %q", reqs[1].Course, api.InformationDesign)
	}
	if reqs[1].RequirementType != api.Optional {
		t.Errorf("Requirements[1].RequirementType: got %q, want %q", reqs[1].RequirementType, api.Optional)
	}
}

func TestSubjectToDetailAPI_EmptySlices(t *testing.T) {
	got := subjectToDetailAPI(domain.Subject{
		ID:                 "subj-1",
		Name:               "Test",
		EligibleAttributes: []domain.SubjectTargetClass{},
		Requirements:       []domain.SubjectRequirement{},
		Faculties:          []domain.SubjectFaculty{},
	})

	if got.EligibleAttributes == nil {
		t.Fatal("EligibleAttributes: got nil, want non-nil empty slice pointer")
	}
	if len(*got.EligibleAttributes) != 0 {
		t.Errorf("EligibleAttributes len: got %d, want 0", len(*got.EligibleAttributes))
	}
	if got.Requirements == nil {
		t.Fatal("Requirements: got nil, want non-nil empty slice pointer")
	}
	if len(*got.Requirements) != 0 {
		t.Errorf("Requirements len: got %d, want 0", len(*got.Requirements))
	}
}

func TestSubjectsToAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   []domain.Subject
		wantLen int
	}{
		{
			name:    "空のスライス",
			input:   []domain.Subject{},
			wantLen: 0,
		},
		{
			name: "複数の科目",
			input: []domain.Subject{
				{ID: "s-1", Name: "科目A"},
				{ID: "s-2", Name: "科目B"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subjectsToAPI(tt.input)
			if len(got) != tt.wantLen {
				t.Fatalf("len: got %d, want %d", len(got), tt.wantLen)
			}
			for i, s := range got {
				if s.Id != tt.input[i].ID {
					t.Errorf("[%d] Id: got %q, want %q", i, s.Id, tt.input[i].ID)
				}
			}
		})
	}
}

func TestBuildSubjectListFilter_AllParams(t *testing.T) {
	year := 2025
	q := "アルゴリズム"
	ids := []string{"s-1", "s-2"}
	grades := []api.DottoFoundationV1Grade{api.B1, api.B2}
	courses := []api.DottoFoundationV1Course{api.InformationSystem}
	classes := []api.DottoFoundationV1Class{api.A, api.B}
	classifications := []api.DottoFoundationV1SubjectClassification{api.Specialized}
	semesters := []api.DottoFoundationV1CourseSemester{api.Q1, api.Q2}
	reqTypes := []api.DottoFoundationV1SubjectRequirementType{api.Required}
	cats := []api.DottoFoundationV1CulturalSubjectCategory{api.Society}

	got := buildSubjectListFilter(api.SubjectsV1ListParams{
		Ids:                       &ids,
		Q:                         &q,
		Grades:                    &grades,
		Courses:                   &courses,
		Classes:                   &classes,
		Classifications:           &classifications,
		Year:                      &year,
		Semesters:                 &semesters,
		RequirementTypes:          &reqTypes,
		CulturalSubjectCategories: &cats,
	})

	if len(got.IDs) != 2 || got.IDs[0] != "s-1" {
		t.Errorf("IDs: got %v", got.IDs)
	}
	if got.Q == nil || *got.Q != q {
		t.Errorf("Q: got %v, want %q", got.Q, q)
	}
	if len(got.Grade) != 2 || got.Grade[0] != domain.GradeB1 {
		t.Errorf("Grade: got %v", got.Grade)
	}
	if len(got.Courses) != 1 || got.Courses[0] != domain.CourseTypeInformationSystem {
		t.Errorf("Courses: got %v", got.Courses)
	}
	if len(got.Class) != 2 || got.Class[0] != domain.ClassA {
		t.Errorf("Class: got %v", got.Class)
	}
	if len(got.Classification) != 1 || got.Classification[0] != domain.SubjectClassificationSpecialized {
		t.Errorf("Classification: got %v", got.Classification)
	}
	if got.Year == nil || *got.Year != 2025 {
		t.Errorf("Year: got %v, want 2025", got.Year)
	}
	if len(got.Semester) != 2 || got.Semester[0] != domain.CourseSemesterQ1 {
		t.Errorf("Semester: got %v", got.Semester)
	}
	if len(got.RequirementType) != 1 || got.RequirementType[0] != domain.SubjectRequirementTypeRequired {
		t.Errorf("RequirementType: got %v", got.RequirementType)
	}
	if len(got.CulturalSubjectCategory) != 1 || got.CulturalSubjectCategory[0] != domain.CulturalSubjectCategorySociety {
		t.Errorf("CulturalSubjectCategory: got %v", got.CulturalSubjectCategory)
	}
}

func TestBuildSubjectListFilter_AllNil(t *testing.T) {
	got := buildSubjectListFilter(api.SubjectsV1ListParams{})

	if got.IDs != nil {
		t.Errorf("IDs: got %v, want nil", got.IDs)
	}
	if got.Q != nil {
		t.Errorf("Q: got %v, want nil", got.Q)
	}
	if got.Grade != nil {
		t.Errorf("Grade: got %v, want nil", got.Grade)
	}
	if got.Courses != nil {
		t.Errorf("Courses: got %v, want nil", got.Courses)
	}
	if got.Class != nil {
		t.Errorf("Class: got %v, want nil", got.Class)
	}
	if got.Classification != nil {
		t.Errorf("Classification: got %v, want nil", got.Classification)
	}
	if got.Year == nil {
		t.Fatal("Year: got nil, want current year")
	}
	if *got.Year != time.Now().Year() {
		t.Errorf("Year: got %d, want %d", *got.Year, time.Now().Year())
	}
	if got.Semester != nil {
		t.Errorf("Semester: got %v, want nil", got.Semester)
	}
	if got.RequirementType != nil {
		t.Errorf("RequirementType: got %v, want nil", got.RequirementType)
	}
	if got.CulturalSubjectCategory != nil {
		t.Errorf("CulturalSubjectCategory: got %v, want nil", got.CulturalSubjectCategory)
	}
}
