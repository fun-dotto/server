package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/proto"
)

type TripHandler struct {
	c crud[model.Trip]
}

func NewTripHandler(repo Repository[model.Trip]) *TripHandler {
	return &TripHandler{c: crud[model.Trip]{
		name: "trip",
		repo: repo,
		keys: func(m model.Trip) []keyField { return []keyField{{"trip_id", m.TripID}} },
	}}
}

// tripDirectionEnum は GTFS の direction_id（0 / 1）との対応。
var tripDirectionEnum = map[adminv1.TripDirection]int{
	adminv1.TripDirection_TRIP_DIRECTION_OUTBOUND: 0,
	adminv1.TripDirection_TRIP_DIRECTION_INBOUND:  1,
}

var tripColumns = map[string][]string{
	"direction": {"direction_id"},
}

type tripInput interface {
	proto.Message
	GetTripId() string
	GetRouteId() string
	GetServiceId() string
	GetDirection() adminv1.TripDirection
}

func tripFromProto(in tripInput, s fieldSet) (model.Trip, error) {
	var v violations
	var m model.Trip
	if s.has("trip_id") {
		m.TripID = requireString(&v, "trip_id", in.GetTripId())
	}
	if s.has("route_id") {
		m.RouteID = requireString(&v, "route_id", in.GetRouteId())
	}
	if s.has("service_id") {
		m.ServiceID = requireString(&v, "service_id", in.GetServiceId())
	}
	if s.has("direction") && present(in, "direction") {
		d, ok := tripDirectionEnum[in.GetDirection()]
		if !ok {
			v.add("direction", "must be a valid enum value other than UNSPECIFIED")
		}
		m.DirectionID = &d
	}
	return m, v.err()
}

func tripToProto(m model.Trip) *adminv1.Trip {
	t := &adminv1.Trip{TripId: m.TripID, RouteId: m.RouteID, ServiceId: m.ServiceID}
	if m.DirectionID != nil {
		for e, d := range tripDirectionEnum {
			if d == *m.DirectionID {
				t.Direction = ptr(e)
			}
		}
	}
	return t
}

func (h *TripHandler) ListTrips(ctx context.Context, req *connect.Request[adminv1.ListTripsRequest]) (*connect.Response[adminv1.ListTripsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListTripsResponse{NextPageToken: next}
	for _, m := range items {
		res.Trips = append(res.Trips, tripToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *TripHandler) GetTrip(ctx context.Context, req *connect.Request[adminv1.GetTripRequest]) (*connect.Response[adminv1.GetTripResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"trip_id", tripID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetTripResponse{Trip: tripToProto(m)}), nil
}

func (h *TripHandler) CreateTrip(ctx context.Context, req *connect.Request[adminv1.CreateTripRequest]) (*connect.Response[adminv1.CreateTripResponse], error) {
	m, err := tripFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateTripResponse{Trip: tripToProto(created)}), nil
}

func (h *TripHandler) UpdateTrip(ctx context.Context, req *connect.Request[adminv1.UpdateTripRequest]) (*connect.Response[adminv1.UpdateTripResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "trip_id")
	if err != nil {
		return nil, err
	}
	m, err := tripFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"trip_id", tripID}}, &m, s.columns(tripColumns))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateTripResponse{Trip: tripToProto(updated)}), nil
}

func (h *TripHandler) DeleteTrip(ctx context.Context, req *connect.Request[adminv1.DeleteTripRequest]) (*connect.Response[adminv1.DeleteTripResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"trip_id", tripID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteTripResponse{}), nil
}
