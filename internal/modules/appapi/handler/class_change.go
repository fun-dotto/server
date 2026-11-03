package handler

import (
	"time"

	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
	"gorm.io/gorm"
)

type classChangeRequest interface {
	GetSubjectIds() []string
	GetFromDate() *date.Date
	GetUntilDate() *date.Date
}

// classChangeQuery は休講・補講・教室変更の一覧に共通する絞り込みを行う。table はテーブル名。
func (h *Handler) classChangeQuery(ctx context.Context, req classChangeRequest, table string) (*gorm.DB, error) {
	var v violations
	subjectIDs := parseUUIDs(req.GetSubjectIds(), "subject_ids", &v)
	query := preload(h.db.WithContext(ctx), subjectPreloads("Subject.")...).
		Joins("JOIN subjects ON subjects.id = " + table + ".subject_id")
	if len(subjectIDs) > 0 {
		query = query.Where(table+".subject_id IN ?", subjectIDs)
	}
	if d := req.GetFromDate(); d != nil {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("from_date", "must be a valid date")
		}
		query = query.Where(table+".date >= ?", t.Format(dateLayout))
	}
	if d := req.GetUntilDate(); d != nil {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("until_date", "must be a valid date")
		}
		query = query.Where(table+".date <= ?", t.Format(dateLayout))
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	return query.Order(table + ".date ASC, " + table + ".period ASC, subjects.syllabus_id ASC"), nil
}

// ListCancelledClasses は休講の一覧を返す。
func (h *Handler) ListCancelledClasses(ctx context.Context, req *connect.Request[appv1.ListCancelledClassesRequest]) (*connect.Response[appv1.ListCancelledClassesResponse], error) {
	query, err := h.classChangeQuery(ctx, req.Msg, "cancelled_classes")
	if err != nil {
		return nil, err
	}
	var records []model.CancelledClass
	if err := query.Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListCancelledClasses", err)
	}
	return connect.NewResponse(&appv1.ListCancelledClassesResponse{
		CancelledClasses: mapSlice(records, func(c model.CancelledClass) *appv1.CancelledClass {
			return &appv1.CancelledClass{
				Id:      c.ID.String(),
				Subject: toProtoSubject(c.Subject),
				Date:    toProtoDate(c.Date),
				Period:  periodMap.toProto[c.Period],
				Comment: ptrOr(c.Comment),
			}
		}),
	}), nil
}

// ListMakeupClasses は補講の一覧を返す。
func (h *Handler) ListMakeupClasses(ctx context.Context, req *connect.Request[appv1.ListMakeupClassesRequest]) (*connect.Response[appv1.ListMakeupClassesResponse], error) {
	query, err := h.classChangeQuery(ctx, req.Msg, "makeup_classes")
	if err != nil {
		return nil, err
	}
	var records []model.MakeupClass
	if err := query.Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListMakeupClasses", err)
	}
	return connect.NewResponse(&appv1.ListMakeupClassesResponse{
		MakeupClasses: mapSlice(records, func(c model.MakeupClass) *appv1.MakeupClass {
			return &appv1.MakeupClass{
				Id:      c.ID.String(),
				Subject: toProtoSubject(c.Subject),
				Date:    toProtoDate(c.Date),
				Period:  periodMap.toProto[c.Period],
				Comment: ptrOr(c.Comment),
			}
		}),
	}), nil
}

// ListRoomChanges は教室変更の一覧を返す。
func (h *Handler) ListRoomChanges(ctx context.Context, req *connect.Request[appv1.ListRoomChangesRequest]) (*connect.Response[appv1.ListRoomChangesResponse], error) {
	query, err := h.classChangeQuery(ctx, req.Msg, "room_changes")
	if err != nil {
		return nil, err
	}
	var records []model.RoomChange
	if err := query.Preload("OriginalRoom").Preload("NewRoom").Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListRoomChanges", err)
	}
	changes := mapSlice(records, func(c model.RoomChange) *appv1.RoomChange {
		return &appv1.RoomChange{
			Id:           c.ID.String(),
			Subject:      toProtoSubject(c.Subject),
			Date:         toProtoDate(c.Date),
			Period:       periodMap.toProto[c.Period],
			OriginalRoom: toProtoRoom(c.OriginalRoom),
			NewRoom:      toProtoRoom(c.NewRoom),
		}
	})
	rooms := make([]*appv1.Room, 0, len(changes)*2)
	for _, c := range changes {
		rooms = append(rooms, c.OriginalRoom, c.NewRoom)
	}
	if err := h.attachRoomFaculties(ctx, rooms...); err != nil {
		return nil, internalError(ctx, "ListRoomChanges", err)
	}
	return connect.NewResponse(&appv1.ListRoomChangesResponse{RoomChanges: changes}), nil
}

// gormScope は個人カレンダー用に休講・補講・教室変更を科目と期間で絞り込む。
type gormScope struct {
	table       string
	subjectIDs  []string
	from, until time.Time
}

func (s *gormScope) apply(db *gorm.DB) *gorm.DB {
	return db.Where(s.table+".subject_id IN ?", s.subjectIDs).
		Where(s.table+".date >= ? AND "+s.table+".date <= ?", s.from.Format(dateLayout), s.until.Format(dateLayout))
}
