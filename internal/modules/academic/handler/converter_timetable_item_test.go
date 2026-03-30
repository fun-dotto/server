package handler

import (
	"testing"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func TestTimetableItemToAPI_Slot(t *testing.T) {
	baseSubject := domain.Subject{ID: "subject-1", Name: "Test"}

	tests := []struct {
		name     string
		slot     *domain.TimetableSlot
		wantNil  bool
		wantDow api.DottoFoundationV1DayOfWeek
		wantPer api.DottoFoundationV1Period
	}{
		{
			name:    "slot is nil",
			slot:    nil,
			wantNil: true,
		},
		{
			name: "both fields empty string",
			slot: &domain.TimetableSlot{
				DayOfWeek: "",
				Period:    "",
			},
			wantNil: true,
		},
		{
			name: "dayOfWeek empty",
			slot: &domain.TimetableSlot{
				DayOfWeek: "",
				Period:    domain.PeriodPeriod1,
			},
			wantNil: true,
		},
		{
			name: "period empty",
			slot: &domain.TimetableSlot{
				DayOfWeek: domain.DayOfWeekMonday,
				Period:    "",
			},
			wantNil: true,
		},
		{
			name: "whitespace only",
			slot: &domain.TimetableSlot{
				DayOfWeek: " \t",
				Period:    "  ",
			},
			wantNil: true,
		},
		{
			name: "valid slot",
			slot: &domain.TimetableSlot{
				DayOfWeek: domain.DayOfWeekMonday,
				Period:    domain.PeriodPeriod1,
			},
			wantNil: false,
			wantDow: api.Monday,
			wantPer: api.Period1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timetableItemToAPI(domain.TimetableItem{
				ID:      "item-1",
				Subject: baseSubject,
				Slot:    tt.slot,
				Rooms:   nil,
			})

			if tt.wantNil {
				if got.Slot != nil {
					t.Fatalf("Slot: got non-nil %+v, want nil", got.Slot)
				}
				return
			}
			if got.Slot == nil {
				t.Fatal("Slot: got nil, want non-nil")
			}
			if got.Slot.DayOfWeek != tt.wantDow {
				t.Errorf("DayOfWeek: got %q, want %q", got.Slot.DayOfWeek, tt.wantDow)
			}
			if got.Slot.Period != tt.wantPer {
				t.Errorf("Period: got %q, want %q", got.Slot.Period, tt.wantPer)
			}
		})
	}
}

func TestToDomainTimetableItemFromRequest_Slot(t *testing.T) {
	subjectID := "subject-1"

	tests := []struct {
		name    string
		slot    *api.DottoFoundationV1TimetableSlot
		wantNil bool
	}{
		{
			name:    "request slot is nil",
			slot:    nil,
			wantNil: true,
		},
		{
			name: "both empty",
			slot: &api.DottoFoundationV1TimetableSlot{
				DayOfWeek: "",
				Period:    "",
			},
			wantNil: true,
		},
		{
			name: "whitespace only",
			slot: &api.DottoFoundationV1TimetableSlot{
				DayOfWeek: "  ",
				Period:    "\t",
			},
			wantNil: true,
		},
		{
			name: "valid",
			slot: &api.DottoFoundationV1TimetableSlot{
				DayOfWeek: api.Tuesday,
				Period:    api.Period2,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDomainTimetableItemFromRequest(api.TimetableItemRequest{
				SubjectId: subjectID,
				Slot:      tt.slot,
				RoomIds:   []string{},
			})

			if tt.wantNil {
				if got.Slot != nil {
					t.Fatalf("Slot: got %+v, want nil", got.Slot)
				}
				return
			}
			if got.Slot == nil {
				t.Fatal("Slot: got nil, want non-nil")
			}
			if got.Slot.DayOfWeek != domain.DayOfWeekTuesday {
				t.Errorf("DayOfWeek: got %q", got.Slot.DayOfWeek)
			}
			if got.Slot.Period != domain.PeriodPeriod2 {
				t.Errorf("Period: got %q", got.Slot.Period)
			}
		})
	}
}
