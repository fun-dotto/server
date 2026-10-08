package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoomReservationHandler struct {
	c crud[model.RoomReservation]
}

func NewRoomReservationHandler(repo Repository[model.RoomReservation]) *RoomReservationHandler {
	return &RoomReservationHandler{c: crud[model.RoomReservation]{
		name: "room_reservation",
		repo: repo,
		keys: func(m model.RoomReservation) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type roomReservationInput interface {
	GetTitle() string
	GetRoomId() string
	GetStartTime() *timestamppb.Timestamp
	GetEndTime() *timestamppb.Timestamp
}

func roomReservationFromProto(in roomReservationInput, s fieldSet) (model.RoomReservation, error) {
	var v violations
	var m model.RoomReservation
	if s.has("title") {
		m.Title = requireString(&v, "title", in.GetTitle())
	}
	if s.has("room_id") {
		m.RoomID = parseUUID(&v, "room_id", in.GetRoomId())
	}
	if s.has("start_time") {
		m.StartTime = requireTimestamp(&v, "start_time", in.GetStartTime())
	}
	if s.has("end_time") {
		m.EndTime = requireTimestamp(&v, "end_time", in.GetEndTime())
	}
	if s.has("start_time") && s.has("end_time") && len(v) == 0 && !m.EndTime.After(m.StartTime) {
		v.add("end_time", "must be after start_time")
	}
	return m, v.err()
}

func roomReservationToProto(m model.RoomReservation) *adminv1.RoomReservation {
	return &adminv1.RoomReservation{
		Id:        m.ID.String(),
		Title:     m.Title,
		RoomId:    m.RoomID.String(),
		StartTime: toTimestamp(m.StartTime),
		EndTime:   toTimestamp(m.EndTime),
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *RoomReservationHandler) ListRoomReservations(ctx context.Context, req *connect.Request[adminv1.ListRoomReservationsRequest]) (*connect.Response[adminv1.ListRoomReservationsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListRoomReservationsResponse{NextPageToken: next}
	for _, m := range items {
		res.RoomReservations = append(res.RoomReservations, roomReservationToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *RoomReservationHandler) GetRoomReservation(ctx context.Context, req *connect.Request[adminv1.GetRoomReservationRequest]) (*connect.Response[adminv1.GetRoomReservationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetRoomReservationResponse{RoomReservation: roomReservationToProto(m)}), nil
}

func (h *RoomReservationHandler) CreateRoomReservation(ctx context.Context, req *connect.Request[adminv1.CreateRoomReservationRequest]) (*connect.Response[adminv1.CreateRoomReservationResponse], error) {
	m, err := roomReservationFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateRoomReservationResponse{RoomReservation: roomReservationToProto(created)}), nil
}

func (h *RoomReservationHandler) UpdateRoomReservation(ctx context.Context, req *connect.Request[adminv1.UpdateRoomReservationRequest]) (*connect.Response[adminv1.UpdateRoomReservationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := roomReservationFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateRoomReservationResponse{RoomReservation: roomReservationToProto(updated)}), nil
}

func (h *RoomReservationHandler) DeleteRoomReservation(ctx context.Context, req *connect.Request[adminv1.DeleteRoomReservationRequest]) (*connect.Response[adminv1.DeleteRoomReservationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteRoomReservationResponse{}), nil
}
