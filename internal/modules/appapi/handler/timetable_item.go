package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

// timetableOrder は曜日・時限・シラバスIDの順に並べる。
const timetableOrder = "CASE timetable_items.day_of_week " +
	"WHEN 'Monday' THEN 1 WHEN 'Tuesday' THEN 2 WHEN 'Wednesday' THEN 3 WHEN 'Thursday' THEN 4 " +
	"WHEN 'Friday' THEN 5 WHEN 'Saturday' THEN 6 WHEN 'Sunday' THEN 7 ELSE 8 END ASC, " +
	"timetable_items.period ASC, (SELECT syllabuses.lesson_id FROM syllabuses WHERE syllabuses.id = subjects.syllabus_id) ASC"

// findTimetableItems は年度・開講時期で時間割を取得する。year が nil の場合は年度で絞り込まない。
func (h *Handler) findTimetableItems(ctx context.Context, year *int, semesters []string) ([]model.TimetableItem, error) {
	query := preload(h.db.WithContext(ctx), append(subjectPreloads("Subject."), "Rooms.Room")...).
		Joins("JOIN subjects ON subjects.id = timetable_items.subject_id").
		Where("subjects.semester IN ?", semesters)
	if year != nil {
		query = query.Where("subjects.year = ?", *year)
	}
	var records []model.TimetableItem
	if err := query.Order(timetableOrder).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// requiredSemesters は必須の semesters を DB の文字列に変換する。
func requiredSemesters(values []appv1.SubjectSemester, v *violations) []string {
	if len(values) == 0 {
		v.add("semesters", "is required")
		return nil
	}
	return semesterMap.fromProtoList(values, "semesters", v)
}

// ListTimetableItems は時間割の一覧を返す。
func (h *Handler) ListTimetableItems(ctx context.Context, req *connect.Request[appv1.ListTimetableItemsRequest]) (*connect.Response[appv1.ListTimetableItemsResponse], error) {
	var v violations
	semesters := requiredSemesters(req.Msg.GetSemesters(), &v)
	if err := v.err(); err != nil {
		return nil, err
	}
	year := h.defaultYear(req.Msg.Year)
	records, err := h.findTimetableItems(ctx, &year, semesters)
	if err != nil {
		return nil, internalError(ctx, "ListTimetableItems", err)
	}
	items := mapSlice(records, toProtoTimetableItem)
	var rooms []*appv1.Room
	for _, item := range items {
		rooms = append(rooms, item.GetRooms()...)
	}
	if err := h.attachRoomFaculties(ctx, rooms...); err != nil {
		return nil, internalError(ctx, "ListTimetableItems", err)
	}
	return connect.NewResponse(&appv1.ListTimetableItemsResponse{TimetableItems: items}), nil
}

func toProtoTimetableItem(item model.TimetableItem) *appv1.TimetableItem {
	p := &appv1.TimetableItem{
		Id:      item.ID.String(),
		Subject: toProtoSubject(item.Subject),
		Rooms:   timetableRooms(item),
	}
	if item.DayOfWeek != nil && item.Period != nil {
		p.Slot = &appv1.TimetableItemSlot{
			DayOfWeek: dayOfWeekMap.toProto[*item.DayOfWeek],
			Period:    periodMap.toProto[*item.Period],
		}
	}
	return p
}

func timetableRooms(item model.TimetableItem) []*appv1.Room {
	rooms := make([]*appv1.Room, 0, len(item.Rooms))
	for _, r := range item.Rooms {
		if r.Room != nil {
			rooms = append(rooms, toProtoRoom(r.Room))
		}
	}
	return rooms
}
