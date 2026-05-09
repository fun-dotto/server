package handler

import (
	"testing"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestFacultyToAPI(t *testing.T) {
	got := facultyToAPI(domain.Faculty{
		ID:    "fac-1",
		Name:  "田中太郎",
		Email: "tanaka@example.com",
	})

	if got.Id != "fac-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "fac-1")
	}
	if got.Name != "田中太郎" {
		t.Errorf("Name: got %q, want %q", got.Name, "田中太郎")
	}
	if got.Email != "tanaka@example.com" {
		t.Errorf("Email: got %q, want %q", got.Email, "tanaka@example.com")
	}
}

func TestFacultiesToAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   []domain.Faculty
		wantLen int
	}{
		{
			name:    "空のスライス",
			input:   []domain.Faculty{},
			wantLen: 0,
		},
		{
			name: "複数の教員",
			input: []domain.Faculty{
				{ID: "fac-1", Name: "教員A", Email: "a@example.com"},
				{ID: "fac-2", Name: "教員B", Email: "b@example.com"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := facultiesToAPI(tt.input)
			if len(got) != tt.wantLen {
				t.Fatalf("len: got %d, want %d", len(got), tt.wantLen)
			}
			for i, f := range got {
				if f.Id != tt.input[i].ID {
					t.Errorf("[%d] Id: got %q, want %q", i, f.Id, tt.input[i].ID)
				}
				if f.Name != tt.input[i].Name {
					t.Errorf("[%d] Name: got %q, want %q", i, f.Name, tt.input[i].Name)
				}
				if f.Email != tt.input[i].Email {
					t.Errorf("[%d] Email: got %q, want %q", i, f.Email, tt.input[i].Email)
				}
			}
		})
	}
}

func TestToDomainFacultyFromRequest(t *testing.T) {
	got := toDomainFacultyFromRequest("fac-1", api.FacultyRequest{
		Name:  "佐藤花子",
		Email: "sato@example.com",
	})

	if got.ID != "fac-1" {
		t.Errorf("ID: got %q, want %q", got.ID, "fac-1")
	}
	if got.Name != "佐藤花子" {
		t.Errorf("Name: got %q, want %q", got.Name, "佐藤花子")
	}
	if got.Email != "sato@example.com" {
		t.Errorf("Email: got %q, want %q", got.Email, "sato@example.com")
	}
}
