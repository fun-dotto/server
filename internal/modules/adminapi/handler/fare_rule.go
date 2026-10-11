package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/proto"
)

type FareRuleHandler struct {
	c crud[model.FareRule]
}

func NewFareRuleHandler(repo Repository[model.FareRule]) *FareRuleHandler {
	return &FareRuleHandler{c: crud[model.FareRule]{
		name: "fare_rule",
		repo: repo,
		keys: func(m model.FareRule) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type fareRuleInput interface {
	proto.Message
	GetFareId() string
	GetRouteId() string
	GetOriginId() string
	GetDestinationId() string
}

func fareRuleFromProto(in fareRuleInput, s fieldSet) (model.FareRule, error) {
	var v violations
	var m model.FareRule
	if s.has("fare_id") {
		m.FareID = requireString(&v, "fare_id", in.GetFareId())
	}
	if s.has("route_id") {
		m.RouteID = optional(in, "route_id", in.GetRouteId())
	}
	if s.has("origin_id") {
		m.OriginID = optional(in, "origin_id", in.GetOriginId())
	}
	if s.has("destination_id") {
		m.DestinationID = optional(in, "destination_id", in.GetDestinationId())
	}
	return m, v.err()
}

func fareRuleToProto(m model.FareRule) *adminv1.FareRule {
	return &adminv1.FareRule{
		Id:            m.ID.String(),
		FareId:        m.FareID,
		RouteId:       m.RouteID,
		OriginId:      m.OriginID,
		DestinationId: m.DestinationID,
	}
}

func (h *FareRuleHandler) ListFareRules(ctx context.Context, req *connect.Request[adminv1.ListFareRulesRequest]) (*connect.Response[adminv1.ListFareRulesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListFareRulesResponse{NextPageToken: next}
	for _, m := range items {
		res.FareRules = append(res.FareRules, fareRuleToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *FareRuleHandler) GetFareRule(ctx context.Context, req *connect.Request[adminv1.GetFareRuleRequest]) (*connect.Response[adminv1.GetFareRuleResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetFareRuleResponse{FareRule: fareRuleToProto(m)}), nil
}

func (h *FareRuleHandler) CreateFareRule(ctx context.Context, req *connect.Request[adminv1.CreateFareRuleRequest]) (*connect.Response[adminv1.CreateFareRuleResponse], error) {
	m, err := fareRuleFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateFareRuleResponse{FareRule: fareRuleToProto(created)}), nil
}

func (h *FareRuleHandler) UpdateFareRule(ctx context.Context, req *connect.Request[adminv1.UpdateFareRuleRequest]) (*connect.Response[adminv1.UpdateFareRuleResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := fareRuleFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateFareRuleResponse{FareRule: fareRuleToProto(updated)}), nil
}

func (h *FareRuleHandler) DeleteFareRule(ctx context.Context, req *connect.Request[adminv1.DeleteFareRuleRequest]) (*connect.Response[adminv1.DeleteFareRuleResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteFareRuleResponse{}), nil
}
