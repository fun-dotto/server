package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/proto"
)

type StopHandler struct {
	c crud[model.Stop]
}

func NewStopHandler(repo Repository[model.Stop]) *StopHandler {
	return &StopHandler{c: crud[model.Stop]{
		name: "stop",
		repo: repo,
		keys: func(m model.Stop) []keyField { return []keyField{{"stop_id", m.StopID}} },
	}}
}

var stopColumns = map[string][]string{
	"stop_location": {"stop_lat", "stop_lon"},
}

type stopInput interface {
	proto.Message
	GetStopId() string
	GetStopName() string
	GetStopLocation() *latlng.LatLng
	GetZoneId() string
}

func stopFromProto(in stopInput, s fieldSet) (model.Stop, error) {
	var v violations
	var m model.Stop
	if s.has("stop_id") {
		m.StopID = requireString(&v, "stop_id", in.GetStopId())
	}
	if s.has("stop_name") {
		m.StopName = requireString(&v, "stop_name", in.GetStopName())
	}
	if s.has("stop_location") {
		loc := in.GetStopLocation()
		switch {
		case loc == nil:
			v.add("stop_location", "must be set")
		case loc.GetLatitude() < -90 || loc.GetLatitude() > 90:
			v.add("stop_location.latitude", "must be in the range [-90.0, +90.0]")
		case loc.GetLongitude() < -180 || loc.GetLongitude() > 180:
			v.add("stop_location.longitude", "must be in the range [-180.0, +180.0]")
		}
		m.StopLat = decimal.NewFromFloat(loc.GetLatitude()).Round(7)
		m.StopLon = decimal.NewFromFloat(loc.GetLongitude()).Round(7)
	}
	if s.has("zone_id") {
		m.ZoneID = optional(in, "zone_id", in.GetZoneId())
	}
	return m, v.err()
}

func stopToProto(m model.Stop) *adminv1.Stop {
	return &adminv1.Stop{
		StopId:   m.StopID,
		StopName: m.StopName,
		StopLocation: &latlng.LatLng{
			Latitude:  m.StopLat.InexactFloat64(),
			Longitude: m.StopLon.InexactFloat64(),
		},
		ZoneId: m.ZoneID,
	}
}

func (h *StopHandler) ListStops(ctx context.Context, req *connect.Request[adminv1.ListStopsRequest]) (*connect.Response[adminv1.ListStopsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListStopsResponse{NextPageToken: next}
	for _, m := range items {
		res.Stops = append(res.Stops, stopToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *StopHandler) GetStop(ctx context.Context, req *connect.Request[adminv1.GetStopRequest]) (*connect.Response[adminv1.GetStopResponse], error) {
	stopID, err := requireKeyString("stop_id", req.Msg.GetStopId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"stop_id", stopID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetStopResponse{Stop: stopToProto(m)}), nil
}

func (h *StopHandler) CreateStop(ctx context.Context, req *connect.Request[adminv1.CreateStopRequest]) (*connect.Response[adminv1.CreateStopResponse], error) {
	m, err := stopFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateStopResponse{Stop: stopToProto(created)}), nil
}

func (h *StopHandler) UpdateStop(ctx context.Context, req *connect.Request[adminv1.UpdateStopRequest]) (*connect.Response[adminv1.UpdateStopResponse], error) {
	stopID, err := requireKeyString("stop_id", req.Msg.GetStopId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "stop_id")
	if err != nil {
		return nil, err
	}
	m, err := stopFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"stop_id", stopID}}, &m, s.columns(stopColumns))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateStopResponse{Stop: stopToProto(updated)}), nil
}

func (h *StopHandler) DeleteStop(ctx context.Context, req *connect.Request[adminv1.DeleteStopRequest]) (*connect.Response[adminv1.DeleteStopResponse], error) {
	stopID, err := requireKeyString("stop_id", req.Msg.GetStopId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"stop_id", stopID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteStopResponse{}), nil
}
