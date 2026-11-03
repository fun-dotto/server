package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
)

// UserHandler は UserService を実装する。
type UserHandler struct {
	appv1connect.UnimplementedUserServiceHandler
	service UserService
}

// NewUserHandler は UserHandler を返す。
func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetUser はログイン中のユーザーを返す。
func (h *UserHandler) GetUser(ctx context.Context, _ *connect.Request[appv1.GetUserRequest]) (*connect.Response[appv1.GetUserResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	u, err := h.service.GetUser(user.ID)
	if err != nil {
		return nil, toConnectError(ctx, "GetUser", err)
	}
	return connect.NewResponse(&appv1.GetUserResponse{User: toProtoUser(*u)}), nil
}

// UpdateUser はログイン中のユーザーを更新する。存在しない場合は作成する。
func (h *UserHandler) UpdateUser(ctx context.Context, req *connect.Request[appv1.UpdateUserRequest]) (*connect.Response[appv1.UpdateUserResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	m := req.Msg

	// update_mask が空の場合は値が入っているフィールドをすべて更新する。
	paths := map[string]bool{"grade": m.Grade != nil, "course": m.Course != nil, "class": m.Class != nil}
	if mask := m.GetUpdateMask(); len(mask.GetPaths()) > 0 {
		for k := range paths {
			paths[k] = false
		}
		var v violations
		for _, p := range mask.GetPaths() {
			if _, ok := paths[p]; !ok {
				v.add("update_mask", "contains an unknown field: "+p)
				continue
			}
			paths[p] = true
		}
		if err := v.err(); err != nil {
			return nil, err
		}
	}

	// 上流の upsert は全フィールドを置き換えるため、更新しないフィールドは現在の値を引き継ぐ。
	current, err := h.service.GetUser(user.ID)
	if err != nil && !errorsIsUserNotFound(err) {
		return nil, toConnectError(ctx, "UpdateUser", err)
	}
	userReq := domain.UserRequest{Email: user.Email}
	if current != nil {
		userReq.Grade, userReq.Course, userReq.Class = current.Grade, current.Course, current.Class
	}

	var v violations
	if paths["grade"] {
		userReq.Grade = optionalFromProto(gradeMap, m.Grade, "grade", &v)
	}
	if paths["course"] {
		userReq.Course = optionalFromProto(courseMap, m.Course, "course", &v)
	}
	if paths["class"] {
		userReq.Class = optionalFromProto(classMap, m.Class, "class", &v)
	}
	if err := v.err(); err != nil {
		return nil, err
	}

	u, err := h.service.UpsertUser(user.ID, userReq)
	if err != nil {
		return nil, toConnectError(ctx, "UpdateUser", err)
	}
	return connect.NewResponse(&appv1.UpdateUserResponse{User: toProtoUser(*u)}), nil
}

// optionalFromProto は optional な enum を domain の値に変換する。未設定・UNSPECIFIED の場合は nil（値を消す）。
func optionalFromProto[D comparable, P comparable](m bimap[D, P], p *P, field string, v *violations) *D {
	if p == nil {
		return nil
	}
	var zero P
	if *p == zero {
		return nil
	}
	d, ok := m.fromProto[*p]
	if !ok {
		v.add(field, "has an unknown value")
		return nil
	}
	return &d
}

func optionalToProto[D comparable, P comparable](m bimap[D, P], d *D) *P {
	if d == nil {
		return nil
	}
	p, ok := m.toProto[*d]
	if !ok {
		return nil
	}
	return &p
}

func toProtoUser(u domain.User) *appv1.User {
	return &appv1.User{
		Grade:  optionalToProto(gradeMap, u.Grade),
		Course: optionalToProto(courseMap, u.Course),
		Class:  optionalToProto(classMap, u.Class),
	}
}

// FcmTokenHandler は FcmTokenService を実装する。
type FcmTokenHandler struct {
	appv1connect.UnimplementedFcmTokenServiceHandler
	service UserService
}

// NewFcmTokenHandler は FcmTokenHandler を返す。
func NewFcmTokenHandler(s UserService) *FcmTokenHandler {
	return &FcmTokenHandler{service: s}
}

// RegisterFcmToken はログイン中のユーザーの FCM トークンを登録する。
func (h *FcmTokenHandler) RegisterFcmToken(ctx context.Context, req *connect.Request[appv1.RegisterFcmTokenRequest]) (*connect.Response[appv1.RegisterFcmTokenResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	if req.Msg.GetToken() == "" {
		return nil, invalidArgument("INVALID_ARGUMENT", "token", "is required")
	}
	token, err := h.service.UpsertFCMToken(user.ID, domain.FCMTokenRequest{Token: req.Msg.GetToken()})
	if err != nil {
		return nil, toConnectError(ctx, "RegisterFcmToken", err)
	}
	return connect.NewResponse(&appv1.RegisterFcmTokenResponse{FcmToken: &appv1.FcmToken{
		Token:     token.Token,
		CreatedAt: toTimestamp(token.CreatedAt),
		UpdatedAt: toTimestamp(token.UpdatedAt),
	}}), nil
}
