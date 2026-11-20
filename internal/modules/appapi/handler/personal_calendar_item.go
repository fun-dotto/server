package handler

import (
	"context"
	"sort"
	"time"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListPersonalCalendarItems はログイン中のユーザーの個人カレンダーアイテムの一覧を返す。
// 時間割、履修情報、休講・補講・教室変更、振替授業日・休日から構築する。
func (h *Handler) ListPersonalCalendarItems(ctx context.Context, req *connect.Request[appv1.ListPersonalCalendarItemsRequest]) (*connect.Response[appv1.ListPersonalCalendarItemsResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var v violations
	if len(req.Msg.GetDates()) == 0 {
		v.add("dates", "is required")
	}
	dates := make([]time.Time, 0, len(req.Msg.GetDates()))
	for _, d := range req.Msg.GetDates() {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("dates", "must contain only valid dates")
			continue
		}
		dates = append(dates, t)
	}
	if err := v.err(); err != nil {
		return nil, err
	}

	in, err := h.loadCalendarInput(ctx, userID, dates)
	if err != nil {
		return nil, internalError(ctx, "ListPersonalCalendarItems", err)
	}
	entries := buildPersonalCalendar(h.calendar, dates, in)

	items := mapSlice(entries, func(e calendarEntry) *appv1.PersonalCalendarItem {
		rooms := make([]*appv1.Room, 0, len(e.rooms))
		for _, r := range e.rooms {
			rooms = append(rooms, toProtoRoom(r))
		}
		return &appv1.PersonalCalendarItem{
			Date:    toProtoDate(e.date),
			Period:  periodMap.toProto[e.period],
			Subject: toProtoSubject(e.subject),
			Rooms:   rooms,
			State:   e.state,
		}
	})
	var rooms []*appv1.Room
	for _, item := range items {
		rooms = append(rooms, item.GetRooms()...)
	}
	if err := h.attachRoomFaculties(ctx, rooms...); err != nil {
		return nil, internalError(ctx, "ListPersonalCalendarItems", err)
	}
	return connect.NewResponse(&appv1.ListPersonalCalendarItemsResponse{PersonalCalendarItems: items}), nil
}

// calendarInput は個人カレンダーの構築に使うデータ。
type calendarInput struct {
	registeredSubjectIDs map[string]struct{}
	timetableItems       []model.TimetableItem
	cancelledClasses     []model.CancelledClass
	makeupClasses        []model.MakeupClass
	roomChanges          []model.RoomChange
}

func (h *Handler) loadCalendarInput(ctx context.Context, userID string, dates []time.Time) (calendarInput, error) {
	in := calendarInput{registeredSubjectIDs: map[string]struct{}{}}
	db := h.db.WithContext(ctx)

	var subjectIDs []string
	if err := db.Model(&model.CourseRegistration{}).Where("user_id = ?", userID).Pluck("subject_id", &subjectIDs).Error; err != nil {
		return in, err
	}
	if len(subjectIDs) == 0 {
		return in, nil
	}
	for _, id := range subjectIDs {
		in.registeredSubjectIDs[id] = struct{}{}
	}

	year, semesters := termsOf(h.calendar.Terms, dates)
	if len(semesters) == 0 {
		return in, nil
	}
	items, err := h.findTimetableItems(ctx, year, semesters)
	if err != nil {
		return in, err
	}
	in.timetableItems = items

	from, until := dateRange(dates)
	scope := func(table string) *gormScope {
		return &gormScope{table: table, subjectIDs: subjectIDs, from: from, until: until}
	}
	if err := scope("cancelled_classes").apply(preload(db, subjectPreloads("Subject.")...)).Find(&in.cancelledClasses).Error; err != nil {
		return in, err
	}
	if err := scope("makeup_classes").apply(preload(db, subjectPreloads("Subject.")...)).Find(&in.makeupClasses).Error; err != nil {
		return in, err
	}
	if err := scope("room_changes").apply(preload(db, subjectPreloads("Subject.")...).Preload("NewRoom")).Find(&in.roomChanges).Error; err != nil {
		return in, err
	}
	return in, nil
}

// calendarEntry は個人カレンダーアイテム 1 件。
type calendarEntry struct {
	date    time.Time
	period  string
	subject *model.Subject
	rooms   []*model.Room
	state   appv1.PersonalCalendarItem_State
}

// buildPersonalCalendar は個人カレンダーアイテムを組み立てる。
// 優先度は 休講 > 補講 > 教室変更 > 通常 とし、同一の日付・時限・科目には 1 件だけ返す。
func buildPersonalCalendar(cal academiccalendar.Calendar, dates []time.Time, in calendarInput) []calendarEntry {
	if len(in.registeredSubjectIDs) == 0 {
		return []calendarEntry{}
	}

	byDay := map[string][]model.TimetableItem{}
	best := map[string]model.TimetableItem{} // 補講の教室に使う時間割
	for _, item := range in.timetableItems {
		if item.DayOfWeek == nil || item.Period == nil {
			continue
		}
		sid := item.SubjectID.String()
		if _, ok := in.registeredSubjectIDs[sid]; !ok {
			continue
		}
		byDay[*item.DayOfWeek] = append(byDay[*item.DayOfWeek], item)

		// 教室のある時間割を優先し、同条件なら曜日・時限が辞書順で小さいものを選ぶ。
		cur, ok := best[sid]
		if !ok {
			best[sid] = item
			continue
		}
		hasRooms, curHasRooms := len(item.Rooms) > 0, len(cur.Rooms) > 0
		if (hasRooms && !curHasRooms) ||
			(hasRooms == curHasRooms && (*item.DayOfWeek < *cur.DayOfWeek ||
				(*item.DayOfWeek == *cur.DayOfWeek && *item.Period < *cur.Period))) {
			best[sid] = item
		}
	}

	type key struct{ date, period, subjectID string }
	cancelled := map[key]struct{}{}
	for _, c := range in.cancelledClasses {
		cancelled[key{c.Date.Format(dateLayout), c.Period, c.SubjectID.String()}] = struct{}{}
	}
	changed := map[key]model.RoomChange{}
	for _, c := range in.roomChanges {
		changed[key{c.Date.Format(dateLayout), c.Period, c.SubjectID.String()}] = c
	}
	makeups := map[key]model.MakeupClass{}
	for _, c := range in.makeupClasses {
		makeups[key{c.Date.Format(dateLayout), c.Period, c.SubjectID.String()}] = c
	}

	result := map[key]calendarEntry{}
	for _, date := range dates {
		ds := date.Format(dateLayout)
		if _, holiday := cal.Holidays[ds]; holiday {
			continue
		}
		w, ok := cal.SubstituteDays[ds]
		if !ok {
			w = date.Weekday()
		}
		for _, item := range byDay[w.String()] {
			k := key{ds, *item.Period, item.SubjectID.String()}
			e := calendarEntry{date: date, period: *item.Period, subject: item.Subject, rooms: itemRooms(item), state: appv1.PersonalCalendarItem_STATE_NORMAL}
			if _, ok := cancelled[k]; ok {
				e.state = appv1.PersonalCalendarItem_STATE_CANCELLED
			} else if rc, ok := changed[k]; ok {
				e.state = appv1.PersonalCalendarItem_STATE_ROOM_CHANGED
				e.rooms = []*model.Room{rc.NewRoom}
			}
			result[k] = e
		}
		for k, mc := range makeups {
			if k.date != ds {
				continue
			}
			if existing, ok := result[k]; ok && existing.state == appv1.PersonalCalendarItem_STATE_CANCELLED {
				continue
			}
			// 補講の教室は DB にないため、該当科目の通常の時間割の教室で代替する（実際の教室と異なる場合がある）。
			rooms := []*model.Room{}
			if item, ok := best[k.subjectID]; ok {
				rooms = uniqueRooms(itemRooms(item))
			}
			result[k] = calendarEntry{date: date, period: mc.Period, subject: mc.Subject, rooms: rooms, state: appv1.PersonalCalendarItem_STATE_MAKEUP}
		}
	}

	entries := make([]calendarEntry, 0, len(result))
	for _, e := range result {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if !a.date.Equal(b.date) {
			return a.date.Before(b.date)
		}
		if a.period != b.period {
			return a.period < b.period
		}
		return subjectID(a.subject) < subjectID(b.subject)
	})
	return entries
}

func subjectID(s *model.Subject) string {
	if s == nil {
		return ""
	}
	return s.ID.String()
}

func itemRooms(item model.TimetableItem) []*model.Room {
	rooms := make([]*model.Room, 0, len(item.Rooms))
	for _, r := range item.Rooms {
		if r.Room != nil {
			rooms = append(rooms, r.Room)
		}
	}
	return rooms
}

func uniqueRooms(rooms []*model.Room) []*model.Room {
	seen := map[string]struct{}{}
	result := make([]*model.Room, 0, len(rooms))
	for _, r := range rooms {
		if _, ok := seen[r.ID.String()]; ok {
			continue
		}
		seen[r.ID.String()] = struct{}{}
		result = append(result, r)
	}
	return result
}

// termsOf は日付が属する年度と開講時期を返す。年度が 1 つに定まらない場合、年度は nil。
func termsOf(terms []academiccalendar.Term, dates []time.Time) (*int, []string) {
	years := map[int]struct{}{}
	semesters := map[string]struct{}{}
	for _, d := range dates {
		ds := d.Format(dateLayout)
		for _, t := range terms {
			if t.Contains(ds) {
				years[t.Year] = struct{}{}
				semesters[t.Semester] = struct{}{}
			}
		}
	}
	var year *int
	if len(years) == 1 {
		for y := range years {
			year = &y
		}
	}
	result := make([]string, 0, len(semesters))
	for s := range semesters {
		result = append(result, s)
	}
	sort.Strings(result)
	return year, result
}

func dateRange(dates []time.Time) (time.Time, time.Time) {
	from, until := dates[0], dates[0]
	for _, d := range dates[1:] {
		if d.Before(from) {
			from = d
		}
		if d.After(until) {
			until = d
		}
	}
	return from, until
}
