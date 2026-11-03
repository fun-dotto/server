package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"google.golang.org/genproto/googleapis/type/date"
)

// classChangeFilter は休講・補講・教室変更の一覧に共通する絞り込み条件。
type classChangeFilter struct {
	subjectIDs *[]string
	from       *time.Time
	until      *time.Time
}

type classChangeRequest interface {
	GetSubjectIds() []string
	GetFromDate() *date.Date
	GetUntilDate() *date.Date
}

func toClassChangeFilter(req classChangeRequest) (classChangeFilter, error) {
	var f classChangeFilter
	var v violations
	if ids := req.GetSubjectIds(); len(ids) > 0 {
		f.subjectIDs = &ids
	}
	if d := req.GetFromDate(); d != nil {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("from_date", "must be a valid date")
		}
		f.from = &t
	}
	if d := req.GetUntilDate(); d != nil {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("until_date", "must be a valid date")
		}
		f.until = &t
	}
	return f, v.err()
}

// CancelledClassHandler は CancelledClassService を実装する。
type CancelledClassHandler struct {
	appv1connect.UnimplementedCancelledClassServiceHandler
	service AcademicService
}

// NewCancelledClassHandler は CancelledClassHandler を返す。
func NewCancelledClassHandler(s AcademicService) *CancelledClassHandler {
	return &CancelledClassHandler{service: s}
}

// ListCancelledClasses は休講の一覧を返す。
func (h *CancelledClassHandler) ListCancelledClasses(ctx context.Context, req *connect.Request[appv1.ListCancelledClassesRequest]) (*connect.Response[appv1.ListCancelledClassesResponse], error) {
	f, err := toClassChangeFilter(req.Msg)
	if err != nil {
		return nil, err
	}
	classes, err := h.service.GetCancelledClasses(domain.CancelledClassQuery{SubjectIDs: f.subjectIDs, From: f.from, Until: f.until})
	if err != nil {
		return nil, toConnectError(ctx, "ListCancelledClasses", err)
	}
	page, next, err := paginate(req.Msg, classes)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListCancelledClassesResponse{
		CancelledClasses: mapSlice(page, func(c domain.CancelledClass) *appv1.CancelledClass {
			return &appv1.CancelledClass{
				Id:      c.ID,
				Subject: toProtoSubject(c.Subject),
				Date:    toProtoDate(c.Date),
				Period:  periodMap.toProto[c.Period],
				Comment: c.Comment,
			}
		}),
		NextPageToken: next,
	}), nil
}

// MakeupClassHandler は MakeupClassService を実装する。
type MakeupClassHandler struct {
	appv1connect.UnimplementedMakeupClassServiceHandler
	service AcademicService
}

// NewMakeupClassHandler は MakeupClassHandler を返す。
func NewMakeupClassHandler(s AcademicService) *MakeupClassHandler {
	return &MakeupClassHandler{service: s}
}

// ListMakeupClasses は補講の一覧を返す。
func (h *MakeupClassHandler) ListMakeupClasses(ctx context.Context, req *connect.Request[appv1.ListMakeupClassesRequest]) (*connect.Response[appv1.ListMakeupClassesResponse], error) {
	f, err := toClassChangeFilter(req.Msg)
	if err != nil {
		return nil, err
	}
	classes, err := h.service.GetMakeupClasses(domain.MakeupClassQuery{SubjectIDs: f.subjectIDs, From: f.from, Until: f.until})
	if err != nil {
		return nil, toConnectError(ctx, "ListMakeupClasses", err)
	}
	page, next, err := paginate(req.Msg, classes)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListMakeupClassesResponse{
		MakeupClasses: mapSlice(page, func(c domain.MakeupClass) *appv1.MakeupClass {
			return &appv1.MakeupClass{
				Id:      c.ID,
				Subject: toProtoSubject(c.Subject),
				Date:    toProtoDate(c.Date),
				Period:  periodMap.toProto[c.Period],
				Comment: c.Comment,
			}
		}),
		NextPageToken: next,
	}), nil
}

// RoomChangeHandler は RoomChangeService を実装する。
type RoomChangeHandler struct {
	appv1connect.UnimplementedRoomChangeServiceHandler
	service AcademicService
}

// NewRoomChangeHandler は RoomChangeHandler を返す。
func NewRoomChangeHandler(s AcademicService) *RoomChangeHandler {
	return &RoomChangeHandler{service: s}
}

// ListRoomChanges は教室変更の一覧を返す。
func (h *RoomChangeHandler) ListRoomChanges(ctx context.Context, req *connect.Request[appv1.ListRoomChangesRequest]) (*connect.Response[appv1.ListRoomChangesResponse], error) {
	f, err := toClassChangeFilter(req.Msg)
	if err != nil {
		return nil, err
	}
	changes, err := h.service.GetRoomChanges(domain.RoomChangeQuery{SubjectIDs: f.subjectIDs, From: f.from, Until: f.until})
	if err != nil {
		return nil, toConnectError(ctx, "ListRoomChanges", err)
	}
	page, next, err := paginate(req.Msg, changes)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListRoomChangesResponse{
		RoomChanges: mapSlice(page, func(c domain.RoomChange) *appv1.RoomChange {
			return &appv1.RoomChange{
				Id:           c.ID,
				Subject:      toProtoSubject(c.Subject),
				Date:         toProtoDate(c.Date),
				Period:       periodMap.toProto[c.Period],
				OriginalRoom: toProtoRoom(c.OriginalRoom),
				NewRoom:      toProtoRoom(c.NewRoom),
			}
		}),
		NextPageToken: next,
	}), nil
}
