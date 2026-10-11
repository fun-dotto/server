package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type RouteHandler struct {
	c crud[model.Route]
}

func NewRouteHandler(repo Repository[model.Route]) *RouteHandler {
	return &RouteHandler{c: crud[model.Route]{
		name: "route",
		repo: repo,
		keys: func(m model.Route) []keyField { return []keyField{{"route_id", m.RouteID}} },
	}}
}

type routeInput interface {
	GetRouteId() string
	GetRouteShortName() string
}

func routeFromProto(in routeInput, s fieldSet) (model.Route, error) {
	var v violations
	var m model.Route
	if s.has("route_id") {
		m.RouteID = requireString(&v, "route_id", in.GetRouteId())
	}
	if s.has("route_short_name") {
		m.RouteShortName = requireString(&v, "route_short_name", in.GetRouteShortName())
	}
	return m, v.err()
}

func routeToProto(m model.Route) *adminv1.Route {
	return &adminv1.Route{RouteId: m.RouteID, RouteShortName: m.RouteShortName}
}

func (h *RouteHandler) ListRoutes(ctx context.Context, req *connect.Request[adminv1.ListRoutesRequest]) (*connect.Response[adminv1.ListRoutesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListRoutesResponse{NextPageToken: next}
	for _, m := range items {
		res.Routes = append(res.Routes, routeToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *RouteHandler) GetRoute(ctx context.Context, req *connect.Request[adminv1.GetRouteRequest]) (*connect.Response[adminv1.GetRouteResponse], error) {
	routeID, err := requireKeyString("route_id", req.Msg.GetRouteId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"route_id", routeID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetRouteResponse{Route: routeToProto(m)}), nil
}

func (h *RouteHandler) CreateRoute(ctx context.Context, req *connect.Request[adminv1.CreateRouteRequest]) (*connect.Response[adminv1.CreateRouteResponse], error) {
	m, err := routeFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateRouteResponse{Route: routeToProto(created)}), nil
}

func (h *RouteHandler) UpdateRoute(ctx context.Context, req *connect.Request[adminv1.UpdateRouteRequest]) (*connect.Response[adminv1.UpdateRouteResponse], error) {
	routeID, err := requireKeyString("route_id", req.Msg.GetRouteId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "route_id")
	if err != nil {
		return nil, err
	}
	m, err := routeFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"route_id", routeID}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateRouteResponse{Route: routeToProto(updated)}), nil
}

func (h *RouteHandler) DeleteRoute(ctx context.Context, req *connect.Request[adminv1.DeleteRouteRequest]) (*connect.Response[adminv1.DeleteRouteResponse], error) {
	routeID, err := requireKeyString("route_id", req.Msg.GetRouteId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"route_id", routeID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteRouteResponse{}), nil
}
