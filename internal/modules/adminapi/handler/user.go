package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/proto"
)

type UserHandler struct {
	c crud[model.User]
}

func NewUserHandler(repo Repository[model.User]) *UserHandler {
	return &UserHandler{c: crud[model.User]{
		name: "user",
		repo: repo,
		keys: func(m model.User) []keyField { return []keyField{{"id", m.ID}} },
	}}
}

type userInput interface {
	proto.Message
	GetId() string
	GetEmail() string
	GetGrade() adminv1.Grade
	GetCourse() adminv1.Course
	GetClass() adminv1.Class
}

func userFromProto(in userInput, s fieldSet) (model.User, error) {
	var v violations
	var m model.User
	if s.has("id") {
		m.ID = requireString(&v, "id", in.GetId())
	}
	if s.has("email") {
		m.Email = requireString(&v, "email", in.GetEmail())
	}
	if s.has("grade") {
		m.Grade = gradeEnum.toOptionalDB(&v, "grade", optional(in, "grade", in.GetGrade()))
	}
	if s.has("course") {
		m.Course = courseEnum.toOptionalDB(&v, "course", optional(in, "course", in.GetCourse()))
	}
	if s.has("class") {
		m.Class = classEnum.toOptionalDB(&v, "class", optional(in, "class", in.GetClass()))
	}
	return m, v.err()
}

func userToProto(m model.User) *adminv1.User {
	return &adminv1.User{
		Id:     m.ID,
		Email:  m.Email,
		Grade:  gradeEnum.fromOptionalDB(m.Grade),
		Course: courseEnum.fromOptionalDB(m.Course),
		Class:  classEnum.fromOptionalDB(m.Class),
	}
}

func (h *UserHandler) ListUsers(ctx context.Context, req *connect.Request[adminv1.ListUsersRequest]) (*connect.Response[adminv1.ListUsersResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListUsersResponse{NextPageToken: next}
	for _, m := range items {
		res.Users = append(res.Users, userToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *connect.Request[adminv1.GetUserRequest]) (*connect.Response[adminv1.GetUserResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetUserResponse{User: userToProto(m)}), nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *connect.Request[adminv1.CreateUserRequest]) (*connect.Response[adminv1.CreateUserResponse], error) {
	m, err := userFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateUserResponse{User: userToProto(created)}), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *connect.Request[adminv1.UpdateUserRequest]) (*connect.Response[adminv1.UpdateUserResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := userFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateUserResponse{User: userToProto(updated)}), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *connect.Request[adminv1.DeleteUserRequest]) (*connect.Response[adminv1.DeleteUserResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteUserResponse{}), nil
}
