package handler

import (
	"testing"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestCourseRegistrationToAPI(t *testing.T) {
	got := courseRegistrationToAPI(domain.CourseRegistration{
		ID:     "cr-1",
		UserID: "user-1",
		Subject: domain.Subject{
			ID:       "subj-1",
			Name:     "情報工学概論",
			Year:     2026,
			Semester: domain.CourseSemesterQ1,
			Credit:   2,
		},
	})

	if got.Id != "cr-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "cr-1")
	}
	if got.UserId != "user-1" {
		t.Errorf("UserId: got %q, want %q", got.UserId, "user-1")
	}
	if got.Subject.Id != "subj-1" {
		t.Errorf("Subject.Id: got %q, want %q", got.Subject.Id, "subj-1")
	}
	if got.Subject.Name != "情報工学概論" {
		t.Errorf("Subject.Name: got %q, want %q", got.Subject.Name, "情報工学概論")
	}
}

func TestCourseRegistrationsToAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   []domain.CourseRegistration
		wantLen int
	}{
		{
			name:    "空のスライス",
			input:   []domain.CourseRegistration{},
			wantLen: 0,
		},
		{
			name: "複数の履修登録",
			input: []domain.CourseRegistration{
				{ID: "cr-1", UserID: "user-1", Subject: domain.Subject{ID: "s-1"}},
				{ID: "cr-2", UserID: "user-1", Subject: domain.Subject{ID: "s-2"}},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := courseRegistrationsToAPI(tt.input)
			if len(got) != tt.wantLen {
				t.Fatalf("len: got %d, want %d", len(got), tt.wantLen)
			}
			for i, cr := range got {
				if cr.Id != tt.input[i].ID {
					t.Errorf("[%d] Id: got %q, want %q", i, cr.Id, tt.input[i].ID)
				}
			}
		})
	}
}

func TestBuildCourseRegistrationListFilter(t *testing.T) {
	year2025 := 2025

	tests := []struct {
		name           string
		params         api.CourseRegistrationsV1ListParams
		wantUserID     string
		wantYear       int
		wantSemesters  []domain.CourseSemester
	}{
		{
			name: "年度とセメスターを指定",
			params: api.CourseRegistrationsV1ListParams{
				UserId:    "user-1",
				Year:      &year2025,
				Semesters: []api.DottoFoundationV1CourseSemester{api.Q1, api.Q2},
			},
			wantUserID:    "user-1",
			wantYear:      2025,
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterQ1, domain.CourseSemesterQ2},
		},
		{
			name: "年度がnilの場合は現在年度をデフォルトにする",
			params: api.CourseRegistrationsV1ListParams{
				UserId:    "user-2",
				Year:      nil,
				Semesters: nil,
			},
			wantUserID:    "user-2",
			wantYear:      time.Now().Year(),
			wantSemesters: []domain.CourseSemester{},
		},
		{
			name: "セメスターが空の場合",
			params: api.CourseRegistrationsV1ListParams{
				UserId:    "user-3",
				Year:      &year2025,
				Semesters: []api.DottoFoundationV1CourseSemester{},
			},
			wantUserID:    "user-3",
			wantYear:      2025,
			wantSemesters: []domain.CourseSemester{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCourseRegistrationListFilter(tt.params)
			if got.UserID != tt.wantUserID {
				t.Errorf("UserID: got %q, want %q", got.UserID, tt.wantUserID)
			}
			if got.Year == nil {
				t.Fatal("Year: got nil, want non-nil")
			}
			if *got.Year != tt.wantYear {
				t.Errorf("Year: got %d, want %d", *got.Year, tt.wantYear)
			}
			if len(got.Semesters) != len(tt.wantSemesters) {
				t.Fatalf("Semesters len: got %d, want %d", len(got.Semesters), len(tt.wantSemesters))
			}
			for i, s := range got.Semesters {
				if s != tt.wantSemesters[i] {
					t.Errorf("Semesters[%d]: got %q, want %q", i, s, tt.wantSemesters[i])
				}
			}
		})
	}
}

func TestToDomainCourseRegistrationFromRequest(t *testing.T) {
	got := toDomainCourseRegistrationFromRequest(api.CourseRegistrationRequest{
		UserId:    "user-1",
		SubjectId: "subj-1",
	})

	if got.UserID != "user-1" {
		t.Errorf("UserID: got %q, want %q", got.UserID, "user-1")
	}
	if got.Subject.ID != "subj-1" {
		t.Errorf("Subject.ID: got %q, want %q", got.Subject.ID, "subj-1")
	}
}
