package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersV1Detail(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		validate func(t *testing.T, resp api.UsersV1DetailResponseObject, err error)
	}{
		{
			name:    "正常にユーザーが取得できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			validate: func(t *testing.T, resp api.UsersV1DetailResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.UsersV1Detail200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.NotNil(t, result.User.Grade)
				assert.Equal(t, api.DottoFoundationV1Grade("B2"), *result.User.Grade)
				require.NotNil(t, result.User.Course)
				assert.Equal(t, api.DottoFoundationV1Course("InformationSystem"), *result.User.Course)
				require.NotNil(t, result.User.Class)
				assert.Equal(t, api.DottoFoundationV1Class("A"), *result.User.Class)
			},
		},
		{
			name:    "ユーザーが見つからない場合404を返す",
			ctx:     ctxWithUserID("nonexistent"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			validate: func(t *testing.T, resp api.UsersV1DetailResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.UsersV1Detail404Response)
				require.True(t, ok, "レスポンスが404レスポンスではありません")
			},
		},
		{
			name:    "userServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.UsersV1DetailResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errUserServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			validate: func(t *testing.T, resp api.UsersV1DetailResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.UsersV1Detail401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepositoryWithError("getUser", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.UsersV1DetailResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.UsersV1Detail(tt.ctx, api.UsersV1DetailRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}

func TestUsersV1Upsert(t *testing.T) {
	grade := api.DottoFoundationV1Grade("B3")
	course := api.DottoFoundationV1Course("InformationDesign")
	class := api.DottoFoundationV1Class("B")

	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		body     *api.UsersV1UpsertJSONRequestBody
		validate func(t *testing.T, resp api.UsersV1UpsertResponseObject, err error)
	}{
		{
			name:    "正常にユーザーを作成・更新できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			body:    &api.UsersV1UpsertJSONRequestBody{Grade: &grade, Course: &course, Class: &class},
			validate: func(t *testing.T, resp api.UsersV1UpsertResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.UsersV1Upsert200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.NotNil(t, result.User.Grade)
				assert.Equal(t, grade, *result.User.Grade)
				require.NotNil(t, result.User.Course)
				assert.Equal(t, course, *result.User.Course)
				require.NotNil(t, result.User.Class)
				assert.Equal(t, class, *result.User.Class)
			},
		},
		{
			name:    "userServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			body:    &api.UsersV1UpsertJSONRequestBody{Grade: &grade},
			validate: func(t *testing.T, resp api.UsersV1UpsertResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errUserServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			body:    &api.UsersV1UpsertJSONRequestBody{Grade: &grade},
			validate: func(t *testing.T, resp api.UsersV1UpsertResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.UsersV1Upsert401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepositoryWithError("upsert", fmt.Errorf("db error"))))),
			body:    &api.UsersV1UpsertJSONRequestBody{Grade: &grade},
			validate: func(t *testing.T, resp api.UsersV1UpsertResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to upsert user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.UsersV1Upsert(tt.ctx, api.UsersV1UpsertRequestObject{
				Body: tt.body,
			})
			tt.validate(t, resp, err)
		})
	}
}

func TestFCMTokenV1Upsert(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		body     *api.FCMTokenV1UpsertJSONRequestBody
		validate func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error)
	}{
		{
			name:    "正常にFCMトークンを作成・更新できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			body:    &api.FCMTokenV1UpsertJSONRequestBody{Token: "fcm-token-1"},
			validate: func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.FCMTokenV1Upsert200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.Equal(t, "fcm-token-1", result.FcmToken.Token)
				assert.WithinDuration(t, time.Now(), result.FcmToken.CreatedAt, time.Minute)
				assert.WithinDuration(t, time.Now(), result.FcmToken.UpdatedAt, time.Minute)
			},
		},
		{
			name:    "userServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			body:    &api.FCMTokenV1UpsertJSONRequestBody{Token: "fcm-token-1"},
			validate: func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errUserServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			body:    &api.FCMTokenV1UpsertJSONRequestBody{Token: "fcm-token-1"},
			validate: func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.FCMTokenV1Upsert401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "リクエストボディがnilの場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepository()))),
			body:    nil,
			validate: func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "request body is required")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithUserService(service.NewUserService(repository.NewMockUserRepositoryWithError("upsertFCMToken", fmt.Errorf("db error"))))),
			body:    &api.FCMTokenV1UpsertJSONRequestBody{Token: "fcm-token-1"},
			validate: func(t *testing.T, resp api.FCMTokenV1UpsertResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to upsert fcm token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.FCMTokenV1Upsert(tt.ctx, api.FCMTokenV1UpsertRequestObject{
				Body: tt.body,
			})
			tt.validate(t, resp, err)
		})
	}
}
