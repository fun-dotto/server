package handler

import (
	"context"
	"regexp"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/proto"
)

type StopTimeHandler struct {
	c crud[model.StopTime]
}

func NewStopTimeHandler(repo Repository[model.StopTime]) *StopTimeHandler {
	return &StopTimeHandler{c: crud[model.StopTime]{
		name: "stop_time",
		repo: repo,
		keys: func(m model.StopTime) []keyField {
			return []keyField{{"trip_id", m.TripID}, {"stop_sequence", m.StopSequence}}
		},
	}}
}

// gtfsTimePattern は GTFS の時刻（HH:MM:SS。24 時以降も表せる）。
var gtfsTimePattern = regexp.MustCompile(`^\d{1,2}:[0-5]\d:[0-5]\d$`)

type stopTimeInput interface {
	proto.Message
	GetTripId() string
	GetStopSequence() int32
	GetStopId() string
	GetArrivalTime() string
	GetDepartureTime() string
	GetStopHeadsign() string
}

func stopTimeFromProto(in stopTimeInput, s fieldSet) (model.StopTime, error) {
	var v violations
	var m model.StopTime
	if s.has("trip_id") {
		m.TripID = requireString(&v, "trip_id", in.GetTripId())
	}
	if s.has("stop_sequence") {
		m.StopSequence = requireNonNegative(&v, "stop_sequence", in.GetStopSequence())
	}
	if s.has("stop_id") {
		m.StopID = requireString(&v, "stop_id", in.GetStopId())
	}
	if s.has("arrival_time") {
		m.ArrivalTime = requireGTFSTime(&v, "arrival_time", in.GetArrivalTime())
	}
	if s.has("departure_time") {
		m.DepartureTime = requireGTFSTime(&v, "departure_time", in.GetDepartureTime())
	}
	if s.has("stop_headsign") {
		m.StopHeadsign = optional(in, "stop_headsign", in.GetStopHeadsign())
	}
	return m, v.err()
}

func requireGTFSTime(v *violations, field, s string) string {
	if !gtfsTimePattern.MatchString(s) {
		v.add(field, "must be in HH:MM:SS format")
	}
	return s
}

func stopTimeToProto(m model.StopTime) *adminv1.StopTime {
	return &adminv1.StopTime{
		TripId:        m.TripID,
		StopSequence:  int32(m.StopSequence),
		StopId:        m.StopID,
		ArrivalTime:   m.ArrivalTime,
		DepartureTime: m.DepartureTime,
		StopHeadsign:  m.StopHeadsign,
	}
}

func (h *StopTimeHandler) ListStopTimes(ctx context.Context, req *connect.Request[adminv1.ListStopTimesRequest]) (*connect.Response[adminv1.ListStopTimesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListStopTimesResponse{NextPageToken: next}
	for _, m := range items {
		res.StopTimes = append(res.StopTimes, stopTimeToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *StopTimeHandler) GetStopTime(ctx context.Context, req *connect.Request[adminv1.GetStopTimeRequest]) (*connect.Response[adminv1.GetStopTimeResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	stopSequence, err := requireKeyInt("stop_sequence", req.Msg.GetStopSequence(), 0)
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"trip_id", tripID}, {"stop_sequence", stopSequence}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetStopTimeResponse{StopTime: stopTimeToProto(m)}), nil
}

func (h *StopTimeHandler) CreateStopTime(ctx context.Context, req *connect.Request[adminv1.CreateStopTimeRequest]) (*connect.Response[adminv1.CreateStopTimeResponse], error) {
	m, err := stopTimeFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateStopTimeResponse{StopTime: stopTimeToProto(created)}), nil
}

func (h *StopTimeHandler) UpdateStopTime(ctx context.Context, req *connect.Request[adminv1.UpdateStopTimeRequest]) (*connect.Response[adminv1.UpdateStopTimeResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	stopSequence, err := requireKeyInt("stop_sequence", req.Msg.GetStopSequence(), 0)
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "trip_id", "stop_sequence")
	if err != nil {
		return nil, err
	}
	m, err := stopTimeFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"trip_id", tripID}, {"stop_sequence", stopSequence}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateStopTimeResponse{StopTime: stopTimeToProto(updated)}), nil
}

func (h *StopTimeHandler) DeleteStopTime(ctx context.Context, req *connect.Request[adminv1.DeleteStopTimeRequest]) (*connect.Response[adminv1.DeleteStopTimeResponse], error) {
	tripID, err := requireKeyString("trip_id", req.Msg.GetTripId())
	if err != nil {
		return nil, err
	}
	stopSequence, err := requireKeyInt("stop_sequence", req.Msg.GetStopSequence(), 0)
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"trip_id", tripID}, {"stop_sequence", stopSequence}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteStopTimeResponse{}), nil
}
