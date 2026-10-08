package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
)

type CancelledClassHandler struct {
	c crud[model.CancelledClass]
}

func NewCancelledClassHandler(repo Repository[model.CancelledClass]) *CancelledClassHandler {
	return &CancelledClassHandler{c: crud[model.CancelledClass]{
		name: "cancelled_class",
		repo: repo,
		keys: func(m model.CancelledClass) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type cancelledClassInput interface {
	proto.Message
	GetSubjectId() string
	GetDate() *date.Date
	GetPeriod() adminv1.Period
	GetComment() string
}

func cancelledClassFromProto(in cancelledClassInput, s fieldSet) (model.CancelledClass, error) {
	var v violations
	var m model.CancelledClass
	if s.has("subject_id") {
		m.SubjectID = parseUUID(&v, "subject_id", in.GetSubjectId())
	}
	if s.has("date") {
		m.Date = requireDate(&v, "date", in.GetDate())
	}
	if s.has("period") {
		m.Period = periodEnum.toDB(&v, "period", in.GetPeriod(), true)
	}
	if s.has("comment") {
		m.Comment = optional(in, "comment", in.GetComment())
	}
	return m, v.err()
}

func cancelledClassToProto(m model.CancelledClass) *adminv1.CancelledClass {
	return &adminv1.CancelledClass{
		Id:        m.ID.String(),
		SubjectId: m.SubjectID.String(),
		Date:      toDate(m.Date),
		Period:    periodEnum.fromDB(m.Period),
		Comment:   m.Comment,
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *CancelledClassHandler) ListCancelledClasses(ctx context.Context, req *connect.Request[adminv1.ListCancelledClassesRequest]) (*connect.Response[adminv1.ListCancelledClassesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListCancelledClassesResponse{NextPageToken: next}
	for _, m := range items {
		res.CancelledClasses = append(res.CancelledClasses, cancelledClassToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *CancelledClassHandler) GetCancelledClass(ctx context.Context, req *connect.Request[adminv1.GetCancelledClassRequest]) (*connect.Response[adminv1.GetCancelledClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetCancelledClassResponse{CancelledClass: cancelledClassToProto(m)}), nil
}

func (h *CancelledClassHandler) CreateCancelledClass(ctx context.Context, req *connect.Request[adminv1.CreateCancelledClassRequest]) (*connect.Response[adminv1.CreateCancelledClassResponse], error) {
	m, err := cancelledClassFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateCancelledClassResponse{CancelledClass: cancelledClassToProto(created)}), nil
}

func (h *CancelledClassHandler) UpdateCancelledClass(ctx context.Context, req *connect.Request[adminv1.UpdateCancelledClassRequest]) (*connect.Response[adminv1.UpdateCancelledClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := cancelledClassFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateCancelledClassResponse{CancelledClass: cancelledClassToProto(updated)}), nil
}

func (h *CancelledClassHandler) DeleteCancelledClass(ctx context.Context, req *connect.Request[adminv1.DeleteCancelledClassRequest]) (*connect.Response[adminv1.DeleteCancelledClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteCancelledClassResponse{}), nil
}
