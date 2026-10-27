package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type RoomHandler struct {
	c crud[model.Room]
}

func NewRoomHandler(repo Repository[model.Room]) *RoomHandler {
	return &RoomHandler{c: crud[model.Room]{
		name: "room",
		repo: repo,
		keys: func(m model.Room) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

var roomFloorEnum = enumMap[adminv1.RoomFloor]{
	adminv1.RoomFloor_ROOM_FLOOR_1:       "Floor1",
	adminv1.RoomFloor_ROOM_FLOOR_2:       "Floor2",
	adminv1.RoomFloor_ROOM_FLOOR_3:       "Floor3",
	adminv1.RoomFloor_ROOM_FLOOR_4:       "Floor4",
	adminv1.RoomFloor_ROOM_FLOOR_5:       "Floor5",
	adminv1.RoomFloor_ROOM_FLOOR_6:       "Floor6",
	adminv1.RoomFloor_ROOM_FLOOR_7:       "Floor7",
	adminv1.RoomFloor_ROOM_FLOOR_VIRTUAL: "Virtual",
}

type roomInput interface {
	GetName() string
	GetFloor() adminv1.RoomFloor
}

func roomFromProto(in roomInput, s fieldSet) (model.Room, error) {
	var v violations
	var m model.Room
	if s.has("name") {
		m.Name = requireString(&v, "name", in.GetName())
	}
	if s.has("floor") {
		m.Floor = roomFloorEnum.toDB(&v, "floor", in.GetFloor(), true)
	}
	return m, v.err()
}

func roomToProto(m model.Room) *adminv1.Room {
	return &adminv1.Room{
		Id:        m.ID.String(),
		Name:      m.Name,
		Floor:     roomFloorEnum.fromDB(m.Floor),
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *RoomHandler) ListRooms(ctx context.Context, req *connect.Request[adminv1.ListRoomsRequest]) (*connect.Response[adminv1.ListRoomsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListRoomsResponse{NextPageToken: next}
	for _, m := range items {
		res.Rooms = append(res.Rooms, roomToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *RoomHandler) GetRoom(ctx context.Context, req *connect.Request[adminv1.GetRoomRequest]) (*connect.Response[adminv1.GetRoomResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetRoomResponse{Room: roomToProto(m)}), nil
}

func (h *RoomHandler) CreateRoom(ctx context.Context, req *connect.Request[adminv1.CreateRoomRequest]) (*connect.Response[adminv1.CreateRoomResponse], error) {
	m, err := roomFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateRoomResponse{Room: roomToProto(created)}), nil
}

func (h *RoomHandler) UpdateRoom(ctx context.Context, req *connect.Request[adminv1.UpdateRoomRequest]) (*connect.Response[adminv1.UpdateRoomResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := roomFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateRoomResponse{Room: roomToProto(updated)}), nil
}

func (h *RoomHandler) DeleteRoom(ctx context.Context, req *connect.Request[adminv1.DeleteRoomRequest]) (*connect.Response[adminv1.DeleteRoomResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteRoomResponse{}), nil
}
