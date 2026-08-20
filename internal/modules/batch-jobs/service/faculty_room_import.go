package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

// ErrFacultyRoomUnmatched は CSV 内に faculties / rooms と照合できない行があったことを表す。
// この場合 INSERT は 1 件も実行しない。
var ErrFacultyRoomUnmatched = errors.New("faculties / rooms に一致しない行があります")

type ImportFacultyRepository interface {
	ListFaculties(ctx context.Context) ([]domain.Faculty, error)
}

type ImportFacultyRoomRepository interface {
	InsertBatch(ctx context.Context, rows []domain.FacultyRoomInsert) error
}

// FacultyRoomImportService は faculties CSV を読み faculty_rooms へ一括 INSERT する。
// email は faculties.email（小文字化）、room_name は rooms.name（NormalizeRoomName）で照合し、
// 未一致が 1 件でもあれば INSERT せずエラーを返す。
type FacultyRoomImportService struct {
	faculties    ImportFacultyRepository
	rooms        ScrapeRoomRepository
	facultyRooms ImportFacultyRoomRepository
}

func NewFacultyRoomImportService(
	faculties ImportFacultyRepository,
	rooms ScrapeRoomRepository,
	facultyRooms ImportFacultyRoomRepository,
) *FacultyRoomImportService {
	return &FacultyRoomImportService{
		faculties:    faculties,
		rooms:        rooms,
		facultyRooms: facultyRooms,
	}
}

// ParseFacultyRoomCSV は faculties CSV（UTF-8、BOM 許容、ヘッダ行必須、必須カラム email/room_name）
// を読み、year を付与した行の一覧を返す。
func ParseFacultyRoomCSV(year int, r io.Reader) ([]domain.FacultyRoomImportRow, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	emailIdx, roomNameIdx := -1, -1
	for i, name := range header {
		switch strings.TrimSpace(strings.TrimPrefix(name, "\uFEFF")) {
		case "email":
			emailIdx = i
		case "room_name":
			roomNameIdx = i
		}
	}
	if emailIdx < 0 || roomNameIdx < 0 {
		return nil, fmt.Errorf("csv に必須カラムがありません（email, room_name が必要）: %v", header)
	}

	var rows []domain.FacultyRoomImportRow
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row: %w", err)
		}
		row := domain.FacultyRoomImportRow{Year: year}
		if emailIdx < len(record) {
			row.Email = strings.TrimSpace(record[emailIdx])
		}
		if roomNameIdx < len(record) {
			row.RoomName = strings.TrimSpace(record[roomNameIdx])
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// Import は CSV 行を faculties / rooms と照合し、全行一致した場合のみ faculty_rooms へ
// 1 トランザクションで一括 INSERT する。room_name 空欄の行はスキップする。
func (s *FacultyRoomImportService) Import(ctx context.Context, rows []domain.FacultyRoomImportRow) (domain.FacultyRoomImportSummary, error) {
	faculties, err := s.faculties.ListFaculties(ctx)
	if err != nil {
		return domain.FacultyRoomImportSummary{}, fmt.Errorf("list faculties: %w", err)
	}
	rooms, err := s.rooms.ListRooms(ctx)
	if err != nil {
		return domain.FacultyRoomImportSummary{}, fmt.Errorf("list rooms: %w", err)
	}

	facultyIDByEmail := make(map[string]string, len(faculties))
	for _, f := range faculties {
		key := strings.ToLower(strings.TrimSpace(f.Email))
		if key == "" {
			continue
		}
		if _, ok := facultyIDByEmail[key]; !ok {
			facultyIDByEmail[key] = f.ID
		}
	}
	roomIDByName := newRoomMatcher(rooms)

	summary := domain.FacultyRoomImportSummary{}
	unmatchedEmails := map[string]struct{}{}
	unmatchedRooms := map[string]struct{}{}
	var inserts []domain.FacultyRoomInsert

	for _, row := range rows {
		if row.RoomName == "" {
			summary.SkippedBlank++
			continue
		}
		facultyID, facultyOK := facultyIDByEmail[strings.ToLower(row.Email)]
		roomID, roomOK := roomIDByName.resolve(row.RoomName)
		if !facultyOK {
			unmatchedEmails[row.Email] = struct{}{}
		}
		if !roomOK {
			unmatchedRooms[row.RoomName] = struct{}{}
		}
		if !facultyOK || !roomOK {
			continue
		}
		inserts = append(inserts, domain.FacultyRoomInsert{
			FacultyID: facultyID,
			RoomID:    roomID,
			Year:      row.Year,
		})
	}

	for email := range unmatchedEmails {
		summary.UnmatchedEmails = append(summary.UnmatchedEmails, email)
	}
	for name := range unmatchedRooms {
		summary.UnmatchedRooms = append(summary.UnmatchedRooms, name)
	}
	sort.Strings(summary.UnmatchedEmails)
	sort.Strings(summary.UnmatchedRooms)

	if len(summary.UnmatchedEmails) > 0 || len(summary.UnmatchedRooms) > 0 {
		slog.Error(
			"faculties / rooms に一致しない行があるため INSERT を実行しません",
			"unmatched_emails", summary.UnmatchedEmails,
			"unmatched_rooms", summary.UnmatchedRooms,
		)
		return summary, ErrFacultyRoomUnmatched
	}

	if err := s.facultyRooms.InsertBatch(ctx, inserts); err != nil {
		return summary, fmt.Errorf("insert faculty rooms: %w", err)
	}
	summary.Inserted = len(inserts)
	return summary, nil
}
