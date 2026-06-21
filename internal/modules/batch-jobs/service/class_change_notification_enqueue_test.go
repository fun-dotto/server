package service

import (
	"testing"
	"time"
)

func TestNotifyWindowRejectsZeroClassDate(t *testing.T) {
	if _, _, err := notifyWindow(time.Time{}); err == nil {
		t.Fatal("notifyWindow() error = nil, want non-nil")
	}
}

func TestNotifyWindowBuildsJSTWindow(t *testing.T) {
	classDate := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)

	notifyAfter, notifyBefore, err := notifyWindow(classDate)
	if err != nil {
		t.Fatalf("notifyWindow() error = %v", err)
	}

	wantAfter := time.Date(2026, 6, 21, 18, 0, 0, 0, jst)
	wantBefore := time.Date(2026, 6, 22, 0, 0, 0, 0, jst)
	if !notifyAfter.Equal(wantAfter) {
		t.Fatalf("notifyAfter = %v, want %v", notifyAfter, wantAfter)
	}
	if !notifyBefore.Equal(wantBefore) {
		t.Fatalf("notifyBefore = %v, want %v", notifyBefore, wantBefore)
	}
}
