package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type FacultyHandler struct {
	c crud[model.Faculty]
}

func NewFacultyHandler(repo Repository[model.Faculty]) *FacultyHandler {
	return &FacultyHandler{c: crud[model.Faculty]{
		name: "faculty",
		repo: repo,
		keys: func(m model.Faculty) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type facultyInput interface {
	GetName() string
	GetEmail() string
}

func facultyFromProto(in facultyInput, s fieldSet) (model.Faculty, error) {
	var v violations
	var m model.Faculty
	if s.has("name") {
		m.Name = requireString(&v, "name", in.GetName())
	}
	if s.has("email") {
		m.Email = requireString(&v, "email", in.GetEmail())
	}
	return m, v.err()
}

func facultyToProto(m model.Faculty) *adminv1.Faculty {
	return &adminv1.Faculty{
		Id:        m.ID.String(),
		Name:      m.Name,
		Email:     m.Email,
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *FacultyHandler) ListFaculties(ctx context.Context, req *connect.Request[adminv1.ListFacultiesRequest]) (*connect.Response[adminv1.ListFacultiesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListFacultiesResponse{NextPageToken: next}
	for _, m := range items {
		res.Faculties = append(res.Faculties, facultyToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *FacultyHandler) GetFaculty(ctx context.Context, req *connect.Request[adminv1.GetFacultyRequest]) (*connect.Response[adminv1.GetFacultyResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetFacultyResponse{Faculty: facultyToProto(m)}), nil
}

func (h *FacultyHandler) CreateFaculty(ctx context.Context, req *connect.Request[adminv1.CreateFacultyRequest]) (*connect.Response[adminv1.CreateFacultyResponse], error) {
	m, err := facultyFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateFacultyResponse{Faculty: facultyToProto(created)}), nil
}

func (h *FacultyHandler) UpdateFaculty(ctx context.Context, req *connect.Request[adminv1.UpdateFacultyRequest]) (*connect.Response[adminv1.UpdateFacultyResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := facultyFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateFacultyResponse{Faculty: facultyToProto(updated)}), nil
}

func (h *FacultyHandler) DeleteFaculty(ctx context.Context, req *connect.Request[adminv1.DeleteFacultyRequest]) (*connect.Response[adminv1.DeleteFacultyResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteFacultyResponse{}), nil
}
