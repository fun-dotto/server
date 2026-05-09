package handler

import (
	"testing"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestBuildFacultyRoomListFilter(t *testing.T) {
	jst := time.FixedZone("JST", 9*3600)
	fixedNow := time.Date(2026, time.April, 17, 12, 0, 0, 0, jst)

	original := nowFunc
	nowFunc = func() time.Time { return fixedNow }
	t.Cleanup(func() { nowFunc = original })

	year2024 := 2024

	tests := []struct {
		name   string
		params api.FacultyRoomsV1ListParams
		want   int
	}{
		{
			name:   "年度指定なしは現在の年度にフォールバックする",
			params: api.FacultyRoomsV1ListParams{Year: nil},
			want:   2026,
		},
		{
			name:   "年度を指定するとそのまま使われる",
			params: api.FacultyRoomsV1ListParams{Year: &year2024},
			want:   2024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFacultyRoomListFilter(tt.params)
			if got.Year == nil {
				t.Fatalf("Year: got nil, want %d", tt.want)
			}
			if *got.Year != tt.want {
				t.Errorf("Year: got %d, want %d", *got.Year, tt.want)
			}
		})
	}
}

func TestFacultyRoomToAPI(t *testing.T) {
	got := facultyRoomToAPI(domain.FacultyRoom{
		ID:      "fr-1",
		Faculty: domain.Faculty{ID: "faculty-1", Name: "山田太郎", Email: "yamada@example.com"},
		Room:    domain.Room{ID: "room-1", Name: "R101", Floor: domain.FloorFloor1},
		Year:    2026,
	})

	if got.Id != "fr-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "fr-1")
	}
	if got.Faculty.Id != "faculty-1" {
		t.Errorf("Faculty.Id: got %q, want %q", got.Faculty.Id, "faculty-1")
	}
	if got.Room.Id != "room-1" {
		t.Errorf("Room.Id: got %q, want %q", got.Room.Id, "room-1")
	}
	if got.Year != 2026 {
		t.Errorf("Year: got %d, want %d", got.Year, 2026)
	}
}

func TestToDomainFacultyRoomFromRequest(t *testing.T) {
	got := toDomainFacultyRoomFromRequest(api.FacultyRoomRequest{
		FacultyId: "faculty-1",
		RoomId:    "room-1",
		Year:      2026,
	})

	if got.Faculty.ID != "faculty-1" {
		t.Errorf("Faculty.ID: got %q, want %q", got.Faculty.ID, "faculty-1")
	}
	if got.Room.ID != "room-1" {
		t.Errorf("Room.ID: got %q, want %q", got.Room.ID, "room-1")
	}
	if got.Year != 2026 {
		t.Errorf("Year: got %d, want %d", got.Year, 2026)
	}
}
