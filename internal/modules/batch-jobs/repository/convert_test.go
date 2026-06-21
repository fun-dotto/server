package repository

import (
	"testing"
	"time"
)

func TestParseDateAcceptsDriverDateTimeString(t *testing.T) {
	got, err := parseDate("2026-06-22T00:00:00Z")
	if err != nil {
		t.Fatalf("parseDate() error = %v", err)
	}
	want := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)

	if !got.Equal(want) {
		t.Fatalf("parseDate() = %v, want %v", got, want)
	}
}

func TestParseDateRejectsInvalidDate(t *testing.T) {
	if _, err := parseDate("not-a-date"); err == nil {
		t.Fatal("parseDate() error = nil, want non-nil")
	}
}
