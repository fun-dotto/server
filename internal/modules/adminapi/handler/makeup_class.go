package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
)

type MakeupClassHandler struct {
	c crud[model.MakeupClass]
}

func NewMakeupClassHandler(repo Repository[model.MakeupClass]) *MakeupClassHandler {
	return &MakeupClassHandler{c: crud[model.MakeupClass]{
		name: "makeup_class",
		repo: repo,
		keys: func(m model.MakeupClass) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type makeupClassInput interface {
	proto.Message
	GetSubjectId() string
	GetDate() *date.Date
	GetPeriod() adminv1.Period
	GetComment() string
}

func makeupClassFromProto(in makeupClassInput, s fieldSet) (model.MakeupClass, error) {
	var v violations
	var m model.MakeupClass
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

func makeupClassToProto(m model.MakeupClass) *adminv1.MakeupClass {
	return &adminv1.MakeupClass{
		Id:        m.ID.String(),
		SubjectId: m.SubjectID.String(),
		Date:      toDate(m.Date),
		Period:    periodEnum.fromDB(m.Period),
		Comment:   m.Comment,
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *MakeupClassHandler) ListMakeupClasses(ctx context.Context, req *connect.Request[adminv1.ListMakeupClassesRequest]) (*connect.Response[adminv1.ListMakeupClassesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListMakeupClassesResponse{NextPageToken: next}
	for _, m := range items {
		res.MakeupClasses = append(res.MakeupClasses, makeupClassToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *MakeupClassHandler) GetMakeupClass(ctx context.Context, req *connect.Request[adminv1.GetMakeupClassRequest]) (*connect.Response[adminv1.GetMakeupClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetMakeupClassResponse{MakeupClass: makeupClassToProto(m)}), nil
}

func (h *MakeupClassHandler) CreateMakeupClass(ctx context.Context, req *connect.Request[adminv1.CreateMakeupClassRequest]) (*connect.Response[adminv1.CreateMakeupClassResponse], error) {
	m, err := makeupClassFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateMakeupClassResponse{MakeupClass: makeupClassToProto(created)}), nil
}

func (h *MakeupClassHandler) UpdateMakeupClass(ctx context.Context, req *connect.Request[adminv1.UpdateMakeupClassRequest]) (*connect.Response[adminv1.UpdateMakeupClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := makeupClassFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateMakeupClassResponse{MakeupClass: makeupClassToProto(updated)}), nil
}

func (h *MakeupClassHandler) DeleteMakeupClass(ctx context.Context, req *connect.Request[adminv1.DeleteMakeupClassRequest]) (*connect.Response[adminv1.DeleteMakeupClassResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteMakeupClassResponse{}), nil
}
