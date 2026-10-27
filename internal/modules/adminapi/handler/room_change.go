package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
)

type RoomChangeHandler struct {
	c crud[model.RoomChange]
}

func NewRoomChangeHandler(repo Repository[model.RoomChange]) *RoomChangeHandler {
	return &RoomChangeHandler{c: crud[model.RoomChange]{
		name: "room_change",
		repo: repo,
		keys: func(m model.RoomChange) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type roomChangeInput interface {
	GetSubjectId() string
	GetDate() *date.Date
	GetPeriod() adminv1.Period
	GetOriginalRoomId() string
	GetNewRoomId() string
}

func roomChangeFromProto(in roomChangeInput, s fieldSet) (model.RoomChange, error) {
	var v violations
	var m model.RoomChange
	if s.has("subject_id") {
		m.SubjectID = parseUUID(&v, "subject_id", in.GetSubjectId())
	}
	if s.has("date") {
		m.Date = requireDate(&v, "date", in.GetDate())
	}
	if s.has("period") {
		m.Period = periodEnum.toDB(&v, "period", in.GetPeriod(), true)
	}
	if s.has("original_room_id") {
		m.OriginalRoomID = parseUUID(&v, "original_room_id", in.GetOriginalRoomId())
	}
	if s.has("new_room_id") {
		m.NewRoomID = parseUUID(&v, "new_room_id", in.GetNewRoomId())
	}
	return m, v.err()
}

func roomChangeToProto(m model.RoomChange) *adminv1.RoomChange {
	return &adminv1.RoomChange{
		Id:             m.ID.String(),
		SubjectId:      m.SubjectID.String(),
		Date:           toDate(m.Date),
		Period:         periodEnum.fromDB(m.Period),
		OriginalRoomId: m.OriginalRoomID.String(),
		NewRoomId:      m.NewRoomID.String(),
		CreatedAt:      toTimestamp(m.CreatedAt),
		UpdatedAt:      toTimestamp(m.UpdatedAt),
	}
}

func (h *RoomChangeHandler) ListRoomChanges(ctx context.Context, req *connect.Request[adminv1.ListRoomChangesRequest]) (*connect.Response[adminv1.ListRoomChangesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListRoomChangesResponse{NextPageToken: next}
	for _, m := range items {
		res.RoomChanges = append(res.RoomChanges, roomChangeToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *RoomChangeHandler) GetRoomChange(ctx context.Context, req *connect.Request[adminv1.GetRoomChangeRequest]) (*connect.Response[adminv1.GetRoomChangeResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetRoomChangeResponse{RoomChange: roomChangeToProto(m)}), nil
}

func (h *RoomChangeHandler) CreateRoomChange(ctx context.Context, req *connect.Request[adminv1.CreateRoomChangeRequest]) (*connect.Response[adminv1.CreateRoomChangeResponse], error) {
	m, err := roomChangeFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateRoomChangeResponse{RoomChange: roomChangeToProto(created)}), nil
}

func (h *RoomChangeHandler) UpdateRoomChange(ctx context.Context, req *connect.Request[adminv1.UpdateRoomChangeRequest]) (*connect.Response[adminv1.UpdateRoomChangeResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := roomChangeFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateRoomChangeResponse{RoomChange: roomChangeToProto(updated)}), nil
}

func (h *RoomChangeHandler) DeleteRoomChange(ctx context.Context, req *connect.Request[adminv1.DeleteRoomChangeRequest]) (*connect.Response[adminv1.DeleteRoomChangeResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteRoomChangeResponse{}), nil
}
