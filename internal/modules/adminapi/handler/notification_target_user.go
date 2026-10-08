package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationTargetUserHandler struct {
	c crud[model.NotificationTargetUser]
}

func NewNotificationTargetUserHandler(repo Repository[model.NotificationTargetUser]) *NotificationTargetUserHandler {
	return &NotificationTargetUserHandler{c: crud[model.NotificationTargetUser]{
		name: "notification_target_user",
		repo: repo,
		keys: func(m model.NotificationTargetUser) []keyField {
			return []keyField{{"notification_id", m.NotificationID.String()}, {"user_id", m.UserID}}
		},
	}}
}

type notificationTargetUserInput interface {
	GetNotificationId() string
	GetUserId() string
	GetNotifiedAt() *timestamppb.Timestamp
}

func notificationTargetUserFromProto(in notificationTargetUserInput, s fieldSet) (model.NotificationTargetUser, error) {
	var v violations
	var m model.NotificationTargetUser
	if s.has("notification_id") {
		m.NotificationID = parseUUID(&v, "notification_id", in.GetNotificationId())
	}
	if s.has("user_id") {
		m.UserID = requireString(&v, "user_id", in.GetUserId())
	}
	if s.has("notified_at") {
		m.NotifiedAt = optionalTimestamp(&v, "notified_at", in.GetNotifiedAt())
	}
	return m, v.err()
}

func notificationTargetUserToProto(m model.NotificationTargetUser) *adminv1.NotificationTargetUser {
	return &adminv1.NotificationTargetUser{
		NotificationId: m.NotificationID.String(),
		UserId:         m.UserID,
		NotifiedAt:     toOptionalTimestamp(m.NotifiedAt),
	}
}

func (h *NotificationTargetUserHandler) ListNotificationTargetUsers(ctx context.Context, req *connect.Request[adminv1.ListNotificationTargetUsersRequest]) (*connect.Response[adminv1.ListNotificationTargetUsersResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListNotificationTargetUsersResponse{NextPageToken: next}
	for _, m := range items {
		res.NotificationTargetUsers = append(res.NotificationTargetUsers, notificationTargetUserToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *NotificationTargetUserHandler) GetNotificationTargetUser(ctx context.Context, req *connect.Request[adminv1.GetNotificationTargetUserRequest]) (*connect.Response[adminv1.GetNotificationTargetUserResponse], error) {
	notificationID, err := parseIDField("notification_id", req.Msg.GetNotificationId())
	if err != nil {
		return nil, err
	}
	userID, err := requireKeyString("user_id", req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"notification_id", notificationID}, {"user_id", userID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetNotificationTargetUserResponse{NotificationTargetUser: notificationTargetUserToProto(m)}), nil
}

func (h *NotificationTargetUserHandler) CreateNotificationTargetUser(ctx context.Context, req *connect.Request[adminv1.CreateNotificationTargetUserRequest]) (*connect.Response[adminv1.CreateNotificationTargetUserResponse], error) {
	m, err := notificationTargetUserFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateNotificationTargetUserResponse{NotificationTargetUser: notificationTargetUserToProto(created)}), nil
}

func (h *NotificationTargetUserHandler) UpdateNotificationTargetUser(ctx context.Context, req *connect.Request[adminv1.UpdateNotificationTargetUserRequest]) (*connect.Response[adminv1.UpdateNotificationTargetUserResponse], error) {
	notificationID, err := parseIDField("notification_id", req.Msg.GetNotificationId())
	if err != nil {
		return nil, err
	}
	userID, err := requireKeyString("user_id", req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "notification_id", "user_id")
	if err != nil {
		return nil, err
	}
	m, err := notificationTargetUserFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"notification_id", notificationID}, {"user_id", userID}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateNotificationTargetUserResponse{NotificationTargetUser: notificationTargetUserToProto(updated)}), nil
}

func (h *NotificationTargetUserHandler) DeleteNotificationTargetUser(ctx context.Context, req *connect.Request[adminv1.DeleteNotificationTargetUserRequest]) (*connect.Response[adminv1.DeleteNotificationTargetUserResponse], error) {
	notificationID, err := parseIDField("notification_id", req.Msg.GetNotificationId())
	if err != nil {
		return nil, err
	}
	userID, err := requireKeyString("user_id", req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"notification_id", notificationID}, {"user_id", userID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteNotificationTargetUserResponse{}), nil
}
