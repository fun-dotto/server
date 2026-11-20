package academiccalendar

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Terms) == 0 || len(c.Holidays) == 0 {
		t.Fatalf("calendar is empty: %+v", c)
	}
}

func TestCurrentAcademicYear(t *testing.T) {
	jst := time.FixedZone("JST", 9*3600)
	tests := []struct {
		now  time.Time
		want int
	}{
		{time.Date(2026, 4, 1, 0, 0, 0, 0, jst), 2026},
		{time.Date(2027, 3, 31, 23, 59, 0, 0, jst), 2026},
		{time.Date(2026, 3, 31, 15, 30, 0, 0, time.UTC), 2026}, // JST では 4 月 1 日
	}
	for _, tt := range tests {
		if got := CurrentAcademicYear(tt.now); got != tt.want {
			t.Errorf("CurrentAcademicYear(%v) = %d, want %d", tt.now, got, tt.want)
		}
	}
}
