package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

type stubFacultyRepository struct {
	faculties []domain.Faculty
}

func (r *stubFacultyRepository) ListFaculties(_ context.Context) ([]domain.Faculty, error) {
	return r.faculties, nil
}

type stubFacultyRoomRepository struct {
	inserted []domain.FacultyRoomInsert
}

func (r *stubFacultyRoomRepository) InsertBatch(_ context.Context, rows []domain.FacultyRoomInsert) error {
	r.inserted = append(r.inserted, rows...)
	return nil
}

func newImportServiceForTest() (*FacultyRoomImportService, *stubFacultyRoomRepository) {
	facultyRooms := &stubFacultyRoomRepository{}
	svc := NewFacultyRoomImportService(
		&stubFacultyRepository{faculties: []domain.Faculty{
			{ID: "faculty-1", Email: "Alice@fun.ac.jp"},
			{ID: "faculty-2", Email: "bob@fun.ac.jp"},
		}},
		&stubRoomRepository{rooms: []domain.RoomRef{
			{ID: "room-1", Name: "R791"},
			{ID: "room-2", Name: "R792"},
		}},
		facultyRooms,
	)
	return svc, facultyRooms
}

func TestParseFacultyRoomCSV(t *testing.T) {
	csvData := "\uFEFFemail,room_name,note\nalice@fun.ac.jp,R791,x\nbob@fun.ac.jp,,y\n"
	rows, err := ParseFacultyRoomCSV(2026, strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseFacultyRoomCSV: %v", err)
	}
	want := []domain.FacultyRoomImportRow{
		{Year: 2026, Email: "alice@fun.ac.jp", RoomName: "R791"},
		{Year: 2026, Email: "bob@fun.ac.jp", RoomName: ""},
	}
	if len(rows) != len(want) {
		t.Fatalf("rows = %d, want %d", len(rows), len(want))
	}
	for i := range want {
		if rows[i] != want[i] {
			t.Errorf("rows[%d] = %+v, want %+v", i, rows[i], want[i])
		}
	}
}

func TestParseFacultyRoomCSVMissingColumn(t *testing.T) {
	if _, err := ParseFacultyRoomCSV(2026, strings.NewReader("email,name\na@fun.ac.jp,R791\n")); err == nil {
		t.Fatal("expected error for missing room_name column")
	}
}

func TestFacultyRoomImport(t *testing.T) {
	svc, facultyRooms := newImportServiceForTest()

	summary, err := svc.Import(context.Background(), []domain.FacultyRoomImportRow{
		{Year: 2026, Email: "alice@fun.ac.jp", RoomName: "R791"},
		{Year: 2026, Email: "BOB@fun.ac.jp", RoomName: "R792"},
		{Year: 2026, Email: "carol@fun.ac.jp", RoomName: ""},
	})
	if err != nil {
		t.Fatalf("Import: %v", err)
	}
	if summary.Inserted != 2 || summary.SkippedBlank != 1 {
		t.Errorf("summary = %+v, want inserted=2 skipped_blank=1", summary)
	}
	if len(facultyRooms.inserted) != 2 {
		t.Fatalf("inserted rows = %d, want 2", len(facultyRooms.inserted))
	}
	if facultyRooms.inserted[0] != (domain.FacultyRoomInsert{FacultyID: "faculty-1", RoomID: "room-1", Year: 2026}) {
		t.Errorf("unexpected first insert: %+v", facultyRooms.inserted[0])
	}
}

func TestFacultyRoomImportAbortsOnUnmatched(t *testing.T) {
	svc, facultyRooms := newImportServiceForTest()

	summary, err := svc.Import(context.Background(), []domain.FacultyRoomImportRow{
		{Year: 2026, Email: "alice@fun.ac.jp", RoomName: "R791"},
		{Year: 2026, Email: "unknown@fun.ac.jp", RoomName: "R999"},
	})
	if !errors.Is(err, ErrFacultyRoomUnmatched) {
		t.Fatalf("err = %v, want ErrFacultyRoomUnmatched", err)
	}
	if len(facultyRooms.inserted) != 0 {
		t.Errorf("inserted rows = %d, want 0 (no insert on unmatched)", len(facultyRooms.inserted))
	}
	if len(summary.UnmatchedEmails) != 1 || summary.UnmatchedEmails[0] != "unknown@fun.ac.jp" {
		t.Errorf("unmatched emails = %v", summary.UnmatchedEmails)
	}
	if len(summary.UnmatchedRooms) != 1 || summary.UnmatchedRooms[0] != "R999" {
		t.Errorf("unmatched rooms = %v", summary.UnmatchedRooms)
	}
}
