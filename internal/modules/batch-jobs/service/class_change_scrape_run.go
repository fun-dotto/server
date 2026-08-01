package service

import (
	"context"
	"fmt"
	"log/slog"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/scraper"
)

// ScrapeSummary はスクレイプ・突合・DB 反映の結果集計。
type ScrapeSummary struct {
	CancelledInserted     int
	CancelledDuplicated   int
	CancelledSkipped      int
	MakeupInserted        int
	MakeupDuplicated      int
	MakeupSkipped         int
	RoomChangeInserted    int
	RoomChangeDuplicated  int
	RoomChangeSkipped     int
	UnmatchedSubjectNames []string
	UnmatchedRoomNames    []string
}

// Run はポータルへログインして休講/補講/部屋変更を取得し、subjects/rooms と突合したうえで
// 冪等に DB へ反映する。year はログイン対象年度（4/1 起点）で、日付の年度補正にも使う。
func (s *ClassChangeScrapeService) Run(ctx context.Context, userID, password string, year int) (ScrapeSummary, error) {
	html, err := s.client.FetchClassChangeHTML(ctx, userID, password)
	if err != nil {
		return ScrapeSummary{}, fmt.Errorf("fetch class change html: %w", err)
	}

	cancelledRows, makeupRows, roomChangeRows, err := scraper.ParseClassChanges(html, year)
	if err != nil {
		return ScrapeSummary{}, fmt.Errorf("parse class change html: %w", err)
	}

	subjectRefs, err := s.subjects.ListSyllabusNames(ctx)
	if err != nil {
		return ScrapeSummary{}, fmt.Errorf("list subjects: %w", err)
	}
	roomRefs, err := s.rooms.ListRooms(ctx)
	if err != nil {
		return ScrapeSummary{}, fmt.Errorf("list rooms: %w", err)
	}

	matcher := newSubjectMatcher(subjectRefs)
	roomMatcher := newRoomMatcher(roomRefs)

	summary := ScrapeSummary{}
	unmatchedSubjects := map[string]struct{}{}
	unmatchedRooms := map[string]struct{}{}

	for _, row := range cancelledRows {
		subjectID, ok := matcher.resolve(row.LessonName)
		if !ok {
			unmatchedSubjects[row.LessonName] = struct{}{}
			summary.CancelledSkipped++
			continue
		}
		if row.Period == "" {
			summary.CancelledSkipped++
			continue
		}

		exists, err := s.cancelled.Exists(ctx, subjectID, row.Date, row.Period)
		if err != nil {
			return ScrapeSummary{}, fmt.Errorf("check cancelled class exists: %w", err)
		}
		if exists {
			summary.CancelledDuplicated++
			continue
		}
		if _, err := s.cancelled.Insert(ctx, academicdomain.CancelledClass{
			Subject: academicdomain.Subject{ID: subjectID},
			Date:    row.Date.Format(dateOutputLayout),
			Period:  academicdomain.Period(row.Period),
			Comment: row.Comment,
		}); err != nil {
			return ScrapeSummary{}, fmt.Errorf("insert cancelled class: %w", err)
		}
		summary.CancelledInserted++
	}

	for _, row := range makeupRows {
		subjectID, ok := matcher.resolve(row.LessonName)
		if !ok {
			unmatchedSubjects[row.LessonName] = struct{}{}
			summary.MakeupSkipped++
			continue
		}
		if row.Period == "" {
			summary.MakeupSkipped++
			continue
		}

		exists, err := s.makeup.Exists(ctx, subjectID, row.Date, row.Period)
		if err != nil {
			return ScrapeSummary{}, fmt.Errorf("check makeup class exists: %w", err)
		}
		if exists {
			summary.MakeupDuplicated++
			continue
		}
		if _, err := s.makeup.Insert(ctx, academicdomain.MakeupClass{
			Subject: academicdomain.Subject{ID: subjectID},
			Date:    row.Date.Format(dateOutputLayout),
			Period:  academicdomain.Period(row.Period),
			Comment: row.Comment,
		}); err != nil {
			return ScrapeSummary{}, fmt.Errorf("insert makeup class: %w", err)
		}
		summary.MakeupInserted++
	}

	for _, row := range roomChangeRows {
		subjectID, subjectOK := matcher.resolve(row.LessonName)
		if !subjectOK {
			unmatchedSubjects[row.LessonName] = struct{}{}
		}
		originalRoomID, fromOK := roomMatcher.resolve(row.RoomFromName)
		if !fromOK {
			unmatchedRooms[row.RoomFromName] = struct{}{}
		}
		newRoomID, toOK := roomMatcher.resolve(row.RoomToName)
		if !toOK {
			unmatchedRooms[row.RoomToName] = struct{}{}
		}
		if !subjectOK || !fromOK || !toOK || row.Period == "" {
			summary.RoomChangeSkipped++
			continue
		}

		exists, err := s.roomChange.Exists(ctx, subjectID, row.Date, row.Period, originalRoomID, newRoomID)
		if err != nil {
			return ScrapeSummary{}, fmt.Errorf("check room change exists: %w", err)
		}
		if exists {
			summary.RoomChangeDuplicated++
			continue
		}
		if _, err := s.roomChange.Insert(ctx, academicdomain.RoomChange{
			Subject:      academicdomain.Subject{ID: subjectID},
			Date:         row.Date.Format(dateOutputLayout),
			Period:       academicdomain.Period(row.Period),
			OriginalRoom: academicdomain.Room{ID: originalRoomID},
			NewRoom:      academicdomain.Room{ID: newRoomID},
		}); err != nil {
			return ScrapeSummary{}, fmt.Errorf("insert room change: %w", err)
		}
		summary.RoomChangeInserted++
	}

	for name := range unmatchedSubjects {
		summary.UnmatchedSubjectNames = append(summary.UnmatchedSubjectNames, name)
	}
	for name := range unmatchedRooms {
		summary.UnmatchedRoomNames = append(summary.UnmatchedRoomNames, name)
	}
	if len(summary.UnmatchedSubjectNames) > 0 {
		slog.Warn("subjects に一致しない授業名があります", "names", summary.UnmatchedSubjectNames)
	}
	if len(summary.UnmatchedRoomNames) > 0 {
		slog.Warn("rooms に一致しない教室名があります", "names", summary.UnmatchedRoomNames)
	}

	return summary, nil
}

const dateOutputLayout = "2006-01-02"
