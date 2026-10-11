package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type ZoneHandler struct {
	c crud[model.Zone]
}

func NewZoneHandler(repo Repository[model.Zone]) *ZoneHandler {
	return &ZoneHandler{c: crud[model.Zone]{
		name: "zone",
		repo: repo,
		keys: func(m model.Zone) []keyField { return []keyField{{"zone_id", m.ZoneID}} },
	}}
}

type zoneInput interface {
	GetZoneId() string
}

func zoneFromProto(in zoneInput, s fieldSet) (model.Zone, error) {
	var v violations
	var m model.Zone
	if s.has("zone_id") {
		m.ZoneID = requireString(&v, "zone_id", in.GetZoneId())
	}
	return m, v.err()
}

func zoneToProto(m model.Zone) *adminv1.Zone {
	return &adminv1.Zone{ZoneId: m.ZoneID}
}

func (h *ZoneHandler) ListZones(ctx context.Context, req *connect.Request[adminv1.ListZonesRequest]) (*connect.Response[adminv1.ListZonesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListZonesResponse{NextPageToken: next}
	for _, m := range items {
		res.Zones = append(res.Zones, zoneToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *ZoneHandler) GetZone(ctx context.Context, req *connect.Request[adminv1.GetZoneRequest]) (*connect.Response[adminv1.GetZoneResponse], error) {
	zoneID, err := requireKeyString("zone_id", req.Msg.GetZoneId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"zone_id", zoneID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetZoneResponse{Zone: zoneToProto(m)}), nil
}

func (h *ZoneHandler) CreateZone(ctx context.Context, req *connect.Request[adminv1.CreateZoneRequest]) (*connect.Response[adminv1.CreateZoneResponse], error) {
	m, err := zoneFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateZoneResponse{Zone: zoneToProto(created)}), nil
}

func (h *ZoneHandler) UpdateZone(ctx context.Context, req *connect.Request[adminv1.UpdateZoneRequest]) (*connect.Response[adminv1.UpdateZoneResponse], error) {
	zoneID, err := requireKeyString("zone_id", req.Msg.GetZoneId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "zone_id")
	if err != nil {
		return nil, err
	}
	m, err := zoneFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"zone_id", zoneID}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateZoneResponse{Zone: zoneToProto(updated)}), nil
}

func (h *ZoneHandler) DeleteZone(ctx context.Context, req *connect.Request[adminv1.DeleteZoneRequest]) (*connect.Response[adminv1.DeleteZoneResponse], error) {
	zoneID, err := requireKeyString("zone_id", req.Msg.GetZoneId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"zone_id", zoneID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteZoneResponse{}), nil
}
