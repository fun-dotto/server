package scraper

import (
	"os"
	"testing"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

func TestParseClassChanges(t *testing.T) {
	html, err := os.ReadFile("testdata/class_changes.html")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	cancelled, makeup, roomChanges, err := ParseClassChanges(string(html), 2026)
	if err != nil {
		t.Fatalf("ParseClassChanges: %v", err)
	}

	if len(cancelled) != 1 {
		t.Fatalf("cancelled: got %d rows, want 1", len(cancelled))
	}
	c := cancelled[0]
	wantDate := time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC)
	if !c.Date.Equal(wantDate) {
		t.Errorf("cancelled date = %v, want %v", c.Date, wantDate)
	}
	if c.Period != "Period3" {
		t.Errorf("cancelled period = %q, want Period3", c.Period)
	}
	if c.LessonName != "情報科学演習" {
		t.Errorf("cancelled lesson name = %q", c.LessonName)
	}
	if c.Kind != domain.CancelledClassKindMakeupScheduled {
		t.Errorf("cancelled kind = %q, want %q", c.Kind, domain.CancelledClassKindMakeupScheduled)
	}

	if len(makeup) != 1 {
		t.Fatalf("makeup: got %d rows, want 1", len(makeup))
	}
	m := makeup[0]
	if m.RoomName != "R781" {
		t.Errorf("makeup room name = %q, want R781", m.RoomName)
	}
	if m.Period != "Period1" {
		t.Errorf("makeup period = %q, want Period1", m.Period)
	}

	if len(roomChanges) != 1 {
		t.Fatalf("room changes: got %d rows, want 1", len(roomChanges))
	}
	rc := roomChanges[0]
	// 03/02 は年度 2026（4/1始まり）の範囲外なので +1 年に補正される。
	wantRCDate := time.Date(2027, time.March, 2, 0, 0, 0, 0, time.UTC)
	if !rc.Date.Equal(wantRCDate) {
		t.Errorf("room change date = %v, want %v", rc.Date, wantRCDate)
	}
	if rc.RoomFromName != "493" || rc.RoomToName != "R781" {
		t.Errorf("room change from/to = %q/%q", rc.RoomFromName, rc.RoomToName)
	}
	if rc.LessonName != "データ構造（旧:データ構造とアルゴリズム）" {
		t.Errorf("room change lesson name = %q", rc.LessonName)
	}
}

func TestFormatPeriod(t *testing.T) {
	tests := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{in: "1時限", want: "Period1", wantOK: true},
		{in: "6", want: "Period6", wantOK: true},
		{in: "", want: "", wantOK: false},
		{in: "時限不明", want: "", wantOK: false},
	}
	for _, tt := range tests {
		got, ok := formatPeriod(tt.in)
		if ok != tt.wantOK || got != tt.want {
			t.Errorf("formatPeriod(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}
