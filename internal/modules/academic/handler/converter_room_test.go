package handler

import (
	"testing"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestRoomToAPI(t *testing.T) {
	got := roomToAPI(domain.Room{
		ID:    "room-1",
		Name:  "R101",
		Floor: domain.FloorFloor3,
	})

	if got.Id != "room-1" {
		t.Errorf("Id: got %q, want %q", got.Id, "room-1")
	}
	if got.Name != "R101" {
		t.Errorf("Name: got %q, want %q", got.Name, "R101")
	}
	if got.Floor != api.Floor3 {
		t.Errorf("Floor: got %q, want %q", got.Floor, api.Floor3)
	}
	if got.Faculty != nil {
		t.Errorf("Faculty: got %+v, want nil", got.Faculty)
	}
	if got.Number != "" {
		t.Errorf("Number: got %q, want empty", got.Number)
	}
}

func TestRoomsToAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   []domain.Room
		wantLen int
	}{
		{
			name:    "空のスライス",
			input:   []domain.Room{},
			wantLen: 0,
		},
		{
			name: "複数の部屋",
			input: []domain.Room{
				{ID: "room-1", Name: "R101", Floor: domain.FloorFloor1},
				{ID: "room-2", Name: "R202", Floor: domain.FloorFloor2},
				{ID: "room-3", Name: "R303", Floor: domain.FloorFloor3},
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roomsToAPI(tt.input)
			if len(got) != tt.wantLen {
				t.Fatalf("len: got %d, want %d", len(got), tt.wantLen)
			}
			for i, r := range got {
				if r.Id != tt.input[i].ID {
					t.Errorf("[%d] Id: got %q, want %q", i, r.Id, tt.input[i].ID)
				}
			}
		})
	}
}

func TestBuildRoomListFilter(t *testing.T) {
	tests := []struct {
		name       string
		params     api.RoomsV1ListParams
		wantFloors []domain.Floor
	}{
		{
			name:       "フロアがnil",
			params:     api.RoomsV1ListParams{Floors: nil},
			wantFloors: nil,
		},
		{
			name: "フロアを指定",
			params: api.RoomsV1ListParams{
				Floors: &[]api.DottoFoundationV1Floor{api.Floor1, api.Floor5},
			},
			wantFloors: []domain.Floor{domain.FloorFloor1, domain.FloorFloor5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRoomListFilter(tt.params)
			if len(got.Floors) != len(tt.wantFloors) {
				t.Fatalf("Floors len: got %d, want %d", len(got.Floors), len(tt.wantFloors))
			}
			for i, f := range got.Floors {
				if f != tt.wantFloors[i] {
					t.Errorf("Floors[%d]: got %q, want %q", i, f, tt.wantFloors[i])
				}
			}
		})
	}
}

func TestToDomainRoomFromRequest(t *testing.T) {
	got := toDomainRoomFromRequest("room-1", api.RoomRequest{
		Name:  "R501",
		Floor: api.Floor5,
	})

	if got.ID != "room-1" {
		t.Errorf("ID: got %q, want %q", got.ID, "room-1")
	}
	if got.Name != "R501" {
		t.Errorf("Name: got %q, want %q", got.Name, "R501")
	}
	if got.Floor != domain.FloorFloor5 {
		t.Errorf("Floor: got %q, want %q", got.Floor, domain.FloorFloor5)
	}
}
