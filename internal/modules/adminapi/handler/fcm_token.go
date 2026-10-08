package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type FcmTokenHandler struct {
	c crud[model.FCMToken]
}

func NewFcmTokenHandler(repo Repository[model.FCMToken]) *FcmTokenHandler {
	return &FcmTokenHandler{c: crud[model.FCMToken]{
		name: "fcm_token",
		repo: repo,
		keys: func(m model.FCMToken) []keyField { return []keyField{{"token", m.Token}} },
	}}
}

type fcmTokenInput interface {
	GetToken() string
	GetUserId() string
}

func fcmTokenFromProto(in fcmTokenInput, s fieldSet) (model.FCMToken, error) {
	var v violations
	var m model.FCMToken
	if s.has("token") {
		m.Token = requireString(&v, "token", in.GetToken())
	}
	if s.has("user_id") {
		m.UserID = requireString(&v, "user_id", in.GetUserId())
	}
	return m, v.err()
}

func fcmTokenToProto(m model.FCMToken) *adminv1.FcmToken {
	return &adminv1.FcmToken{
		Token:     m.Token,
		UserId:    m.UserID,
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *FcmTokenHandler) ListFcmTokens(ctx context.Context, req *connect.Request[adminv1.ListFcmTokensRequest]) (*connect.Response[adminv1.ListFcmTokensResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListFcmTokensResponse{NextPageToken: next}
	for _, m := range items {
		res.FcmTokens = append(res.FcmTokens, fcmTokenToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *FcmTokenHandler) GetFcmToken(ctx context.Context, req *connect.Request[adminv1.GetFcmTokenRequest]) (*connect.Response[adminv1.GetFcmTokenResponse], error) {
	token, err := requireKeyString("token", req.Msg.GetToken())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"token", token}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetFcmTokenResponse{FcmToken: fcmTokenToProto(m)}), nil
}

func (h *FcmTokenHandler) CreateFcmToken(ctx context.Context, req *connect.Request[adminv1.CreateFcmTokenRequest]) (*connect.Response[adminv1.CreateFcmTokenResponse], error) {
	m, err := fcmTokenFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateFcmTokenResponse{FcmToken: fcmTokenToProto(created)}), nil
}

func (h *FcmTokenHandler) UpdateFcmToken(ctx context.Context, req *connect.Request[adminv1.UpdateFcmTokenRequest]) (*connect.Response[adminv1.UpdateFcmTokenResponse], error) {
	token, err := requireKeyString("token", req.Msg.GetToken())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "token")
	if err != nil {
		return nil, err
	}
	m, err := fcmTokenFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"token", token}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateFcmTokenResponse{FcmToken: fcmTokenToProto(updated)}), nil
}

func (h *FcmTokenHandler) DeleteFcmToken(ctx context.Context, req *connect.Request[adminv1.DeleteFcmTokenRequest]) (*connect.Response[adminv1.DeleteFcmTokenResponse], error) {
	token, err := requireKeyString("token", req.Msg.GetToken())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"token", token}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteFcmTokenResponse{}), nil
}
