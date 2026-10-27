package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/protobuf/proto"
)

type TimetableItemHandler struct {
	c crud[model.TimetableItem]
}

func NewTimetableItemHandler(repo Repository[model.TimetableItem]) *TimetableItemHandler {
	return &TimetableItemHandler{c: crud[model.TimetableItem]{
		name: "timetable_item",
		repo: repo,
		keys: func(m model.TimetableItem) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

// timetableItemColumns は子テーブルとして保存するフィールド。repository.TimetableItemRepository が置き換える。
var timetableItemColumns = map[string][]string{
	"room_ids": {"Rooms"},
}

type timetableItemInput interface {
	proto.Message
	GetSubjectId() string
	GetDayOfWeek() dayofweek.DayOfWeek
	GetPeriod() adminv1.Period
	GetRoomIds() []string
}

func timetableItemFromProto(in timetableItemInput, s fieldSet) (model.TimetableItem, error) {
	var v violations
	var m model.TimetableItem
	if s.has("subject_id") {
		m.SubjectID = parseUUID(&v, "subject_id", in.GetSubjectId())
	}
	if s.has("day_of_week") {
		m.DayOfWeek = dayOfWeekEnum.toOptionalDB(&v, "day_of_week", optional(in, "day_of_week", in.GetDayOfWeek()))
	}
	if s.has("period") {
		m.Period = periodEnum.toOptionalDB(&v, "period", optional(in, "period", in.GetPeriod()))
	}
	if s.has("room_ids") {
		for i, id := range in.GetRoomIds() {
			m.Rooms = append(m.Rooms, model.TimetableItemRoom{
				RoomID: parseUUID(&v, fmt.Sprintf("room_ids.%d", i), id).String(),
			})
		}
	}
	return m, v.err()
}

func timetableItemToProto(m model.TimetableItem) *adminv1.TimetableItem {
	t := &adminv1.TimetableItem{
		Id:        m.ID.String(),
		SubjectId: m.SubjectID.String(),
		DayOfWeek: dayOfWeekEnum.fromOptionalDB(m.DayOfWeek),
		Period:    periodEnum.fromOptionalDB(m.Period),
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
	for _, r := range m.Rooms {
		t.RoomIds = append(t.RoomIds, r.RoomID)
	}
	return t
}

func (h *TimetableItemHandler) ListTimetableItems(ctx context.Context, req *connect.Request[adminv1.ListTimetableItemsRequest]) (*connect.Response[adminv1.ListTimetableItemsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListTimetableItemsResponse{NextPageToken: next}
	for _, m := range items {
		res.TimetableItems = append(res.TimetableItems, timetableItemToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *TimetableItemHandler) GetTimetableItem(ctx context.Context, req *connect.Request[adminv1.GetTimetableItemRequest]) (*connect.Response[adminv1.GetTimetableItemResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetTimetableItemResponse{TimetableItem: timetableItemToProto(m)}), nil
}

func (h *TimetableItemHandler) CreateTimetableItem(ctx context.Context, req *connect.Request[adminv1.CreateTimetableItemRequest]) (*connect.Response[adminv1.CreateTimetableItemResponse], error) {
	m, err := timetableItemFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateTimetableItemResponse{TimetableItem: timetableItemToProto(created)}), nil
}

func (h *TimetableItemHandler) UpdateTimetableItem(ctx context.Context, req *connect.Request[adminv1.UpdateTimetableItemRequest]) (*connect.Response[adminv1.UpdateTimetableItemResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := timetableItemFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(timetableItemColumns))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateTimetableItemResponse{TimetableItem: timetableItemToProto(updated)}), nil
}

func (h *TimetableItemHandler) DeleteTimetableItem(ctx context.Context, req *connect.Request[adminv1.DeleteTimetableItemRequest]) (*connect.Response[adminv1.DeleteTimetableItemResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteTimetableItemResponse{}), nil
}
