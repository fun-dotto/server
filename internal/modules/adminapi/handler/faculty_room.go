package handler

import (
	"context"
	"uuid"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type FacultyRoomHandler struct {
	c crud[model.FacultyRoom]
}

func NewFacultyRoomHandler(repo Repository[model.FacultyRoom]) *FacultyRoomHandler {
	return &FacultyRoomHandler{c: crud[model.FacultyRoom]{
		name: "faculty_room",
		repo: repo,
		keys: func(m model.FacultyRoom) []keyField {
			return []keyField{{"faculty_id", m.FacultyID.String()}, {"room_id", m.RoomID.String()}, {"year", m.Year}}
		},
	}}
}

type facultyRoomInput interface {
	GetFacultyId() string
	GetRoomId() string
	GetYear() int32
}

func facultyRoomFromProto(in facultyRoomInput, s fieldSet) (model.FacultyRoom, error) {
	var v violations
	var m model.FacultyRoom
	if s.has("faculty_id") {
		m.FacultyID = parseUUID(&v, "faculty_id", in.GetFacultyId())
	}
	if s.has("room_id") {
		m.RoomID = parseUUID(&v, "room_id", in.GetRoomId())
	}
	if s.has("year") {
		m.Year = requirePositive(&v, "year", in.GetYear())
	}
	return m, v.err()
}

func facultyRoomToProto(m model.FacultyRoom) *adminv1.FacultyRoom {
	return &adminv1.FacultyRoom{
		FacultyId: m.FacultyID.String(),
		RoomId:    m.RoomID.String(),
		Year:      int32(m.Year),
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *FacultyRoomHandler) ListFacultyRooms(ctx context.Context, req *connect.Request[adminv1.ListFacultyRoomsRequest]) (*connect.Response[adminv1.ListFacultyRoomsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListFacultyRoomsResponse{NextPageToken: next}
	for _, m := range items {
		res.FacultyRooms = append(res.FacultyRooms, facultyRoomToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *FacultyRoomHandler) GetFacultyRoom(ctx context.Context, req *connect.Request[adminv1.GetFacultyRoomRequest]) (*connect.Response[adminv1.GetFacultyRoomResponse], error) {
	facultyID, err := parseIDField("faculty_id", req.Msg.GetFacultyId())
	if err != nil {
		return nil, err
	}
	roomID, err := parseIDField("room_id", req.Msg.GetRoomId())
	if err != nil {
		return nil, err
	}
	year, err := requireKeyInt("year", req.Msg.GetYear(), 1)
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"faculty_id", facultyID}, {"room_id", roomID}, {"year", year}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetFacultyRoomResponse{FacultyRoom: facultyRoomToProto(m)}), nil
}

func (h *FacultyRoomHandler) CreateFacultyRoom(ctx context.Context, req *connect.Request[adminv1.CreateFacultyRoomRequest]) (*connect.Response[adminv1.CreateFacultyRoomResponse], error) {
	m, err := facultyRoomFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	m.Id = uuid.New()
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateFacultyRoomResponse{FacultyRoom: facultyRoomToProto(created)}), nil
}

func (h *FacultyRoomHandler) DeleteFacultyRoom(ctx context.Context, req *connect.Request[adminv1.DeleteFacultyRoomRequest]) (*connect.Response[adminv1.DeleteFacultyRoomResponse], error) {
	facultyID, err := parseIDField("faculty_id", req.Msg.GetFacultyId())
	if err != nil {
		return nil, err
	}
	roomID, err := parseIDField("room_id", req.Msg.GetRoomId())
	if err != nil {
		return nil, err
	}
	year, err := requireKeyInt("year", req.Msg.GetYear(), 1)
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"faculty_id", facultyID}, {"room_id", roomID}, {"year", year}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteFacultyRoomResponse{}), nil
}
