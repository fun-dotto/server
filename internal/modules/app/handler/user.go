package handler

import (
	"context"
	"errors"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/middleware"
)

// UsersV1Detail ユーザーを取得する
func (h *Handler) UsersV1Detail(ctx context.Context, request api.UsersV1DetailRequestObject) (api.UsersV1DetailResponseObject, error) {
	if h.userService == nil {
		return nil, errUserServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return api.UsersV1Detail401Response{}, nil
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return api.UsersV1Detail404Response{}, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return api.UsersV1Detail200JSONResponse{
		User: toApiUserInfo(*user),
	}, nil
}

// UsersV1Upsert ユーザーを作成または更新する
func (h *Handler) UsersV1Upsert(ctx context.Context, request api.UsersV1UpsertRequestObject) (api.UsersV1UpsertResponseObject, error) {
	if h.userService == nil {
		return nil, errUserServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return api.UsersV1Upsert401Response{}, nil
	}

	email, _ := middleware.UserEmailFromContext(ctx)

	req := toDomainUserRequest(*request.Body)
	req.Email = email
	user, err := h.userService.UpsertUser(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return api.UsersV1Upsert200JSONResponse{
		User: toApiUserInfo(*user),
	}, nil
}

// FCMTokenV1Upsert FCMトークンを作成または更新する
func (h *Handler) FCMTokenV1Upsert(ctx context.Context, request api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	if h.userService == nil {
		return nil, errUserServiceNotConfigured
	}
	if request.Body == nil {
		return nil, fmt.Errorf("request body is required")
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return api.FCMTokenV1Upsert401Response{}, nil
	}

	token, err := h.userService.UpsertFCMToken(userID, domain.FCMTokenRequest{Token: request.Body.Token})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert fcm token: %w", err)
	}

	return api.FCMTokenV1Upsert200JSONResponse{
		FcmToken: toApiFCMToken(*token),
	}, nil
}

func toApiUserInfo(user domain.User) api.UserInfo {
	info := api.UserInfo{}
	if user.Grade != nil {
		g := api.DottoFoundationV1Grade(*user.Grade)
		info.Grade = &g
	}
	if user.Course != nil {
		c := api.DottoFoundationV1Course(*user.Course)
		info.Course = &c
	}
	if user.Class != nil {
		cl := api.DottoFoundationV1Class(*user.Class)
		info.Class = &cl
	}
	return info
}

func toDomainUserRequest(body api.UsersV1UpsertJSONRequestBody) domain.UserRequest {
	req := domain.UserRequest{}
	if body.Grade != nil {
		g := domain.Grade(*body.Grade)
		req.Grade = &g
	}
	if body.Course != nil {
		c := domain.Course(*body.Course)
		req.Course = &c
	}
	if body.Class != nil {
		cl := domain.Class(*body.Class)
		req.Class = &cl
	}
	return req
}

func toApiFCMToken(token domain.FCMToken) api.FCMToken {
	return api.FCMToken{
		Token:     token.Token,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}
