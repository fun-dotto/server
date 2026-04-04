package handler

import (
	"testing"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestPersonalCalendarItemToAPI(t *testing.T) {
	date := time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)
	got := personalCalendarItemToAPI(domain.PersonalCalendarItem{
		Date:   date,
		Period: domain.PeriodPeriod3,
		Subject: domain.Subject{
			ID:       "subj-1",
			Name:     "データ構造",
			Year:     2026,
			Semester: domain.CourseSemesterQ1,
			Credit:   2,
		},
		Rooms: []domain.Room{
			{ID: "room-1", Name: "R101", Floor: domain.FloorFloor1},
		},
		Status: domain.PersonalCalendarItemStatusNormal,
	})

	if !got.Date.Equal(date) {
		t.Errorf("Date: got %v, want %v", got.Date, date)
	}
	if got.Period != api.Period3 {
		t.Errorf("Period: got %q, want %q", got.Period, api.Period3)
	}
	if got.Subject.Id != "subj-1" {
		t.Errorf("Subject.Id: got %q, want %q", got.Subject.Id, "subj-1")
	}
	if got.Subject.Name != "データ構造" {
		t.Errorf("Subject.Name: got %q, want %q", got.Subject.Name, "データ構造")
	}
	if len(got.Rooms) != 1 {
		t.Fatalf("Rooms len: got %d, want 1", len(got.Rooms))
	}
	if got.Rooms[0].Id != "room-1" {
		t.Errorf("Rooms[0].Id: got %q, want %q", got.Rooms[0].Id, "room-1")
	}
	if got.Status != api.Normal {
		t.Errorf("Status: got %q, want %q", got.Status, api.Normal)
	}
}

func TestPersonalCalendarItemToAPI_EmptyRooms(t *testing.T) {
	got := personalCalendarItemToAPI(domain.PersonalCalendarItem{
		Date:    time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
		Period:  domain.PeriodPeriod1,
		Subject: domain.Subject{ID: "subj-1"},
		Rooms:   []domain.Room{},
		Status:  domain.PersonalCalendarItemStatusCancelled,
	})

	if len(got.Rooms) != 0 {
		t.Errorf("Rooms len: got %d, want 0", len(got.Rooms))
	}
	if got.Status != api.Cancelled {
		t.Errorf("Status: got %q, want %q", got.Status, api.Cancelled)
	}
}

func TestPersonalCalendarItemToAPI_AllStatuses(t *testing.T) {
	tests := []struct {
		name       string
		status     domain.PersonalCalendarItemStatus
		wantStatus api.DottoFoundationV1PersonalCalendarItemStatus
	}{
		{"通常", domain.PersonalCalendarItemStatusNormal, api.Normal},
		{"休講", domain.PersonalCalendarItemStatusCancelled, api.Cancelled},
		{"補講", domain.PersonalCalendarItemStatusMakeup, api.Makeup},
		{"教室変更", domain.PersonalCalendarItemStatusRoomChanged, api.RoomChanged},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := personalCalendarItemToAPI(domain.PersonalCalendarItem{
				Date:    time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
				Period:  domain.PeriodPeriod1,
				Subject: domain.Subject{ID: "subj-1"},
				Rooms:   []domain.Room{},
				Status:  tt.status,
			})
			if got.Status != tt.wantStatus {
				t.Errorf("Status: got %q, want %q", got.Status, tt.wantStatus)
			}
		})
	}
}

func TestPersonalCalendarItemsToAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   []domain.PersonalCalendarItem
		wantLen int
	}{
		{
			name:    "空のスライス",
			input:   []domain.PersonalCalendarItem{},
			wantLen: 0,
		},
		{
			name: "複数のカレンダーアイテム",
			input: []domain.PersonalCalendarItem{
				{
					Date:    time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
					Period:  domain.PeriodPeriod1,
					Subject: domain.Subject{ID: "s-1"},
					Rooms:   []domain.Room{},
					Status:  domain.PersonalCalendarItemStatusNormal,
				},
				{
					Date:    time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
					Period:  domain.PeriodPeriod2,
					Subject: domain.Subject{ID: "s-2"},
					Rooms:   []domain.Room{},
					Status:  domain.PersonalCalendarItemStatusCancelled,
				},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := personalCalendarItemsToAPI(tt.input)
			if len(got) != tt.wantLen {
				t.Fatalf("len: got %d, want %d", len(got), tt.wantLen)
			}
			for i, item := range got {
				if item.Subject.Id != tt.input[i].Subject.ID {
					t.Errorf("[%d] Subject.Id: got %q, want %q", i, item.Subject.Id, tt.input[i].Subject.ID)
				}
			}
		})
	}
}
