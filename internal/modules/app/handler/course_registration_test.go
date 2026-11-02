package handler

import (
	"context"
	"fmt"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ctxWithUserID(userID string) context.Context {
	return context.WithValue(context.Background(), "userID", userID)
}

func TestCourseRegistrationsV1List(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		validate func(t *testing.T, resp api.CourseRegistrationsV1ListResponseObject, err error)
	}{
		{
			name:    "正常に履修登録一覧が取得できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.CourseRegistrationsV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.Len(t, result.CourseRegistrations, 1)
				assert.Equal(t, "cr1", result.CourseRegistrations[0].Id)
				assert.Equal(t, "s1", result.CourseRegistrations[0].Subject.Id)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1ListResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.CourseRegistrationsV1List401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getRegistrations", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get course registrations")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.CourseRegistrationsV1List(tt.ctx, api.CourseRegistrationsV1ListRequestObject{
				Params: api.CourseRegistrationsV1ListParams{
					Semesters: []api.DottoFoundationV1CourseSemester{api.DottoFoundationV1CourseSemester("Q1")},
				},
			})
			tt.validate(t, resp, err)
		})
	}
}

func TestCourseRegistrationsV1Create(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		validate func(t *testing.T, resp api.CourseRegistrationsV1CreateResponseObject, err error)
	}{
		{
			name:    "正常に履修登録を作成できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1CreateResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.CourseRegistrationsV1Create201JSONResponse)
				require.True(t, ok, "レスポンスが201 JSONレスポンスではありません")
				assert.Equal(t, "cr-new", result.CourseRegistration.Id)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1CreateResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1CreateResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.CourseRegistrationsV1Create401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("create", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1CreateResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create course registration")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.CourseRegistrationsV1Create(tt.ctx, api.CourseRegistrationsV1CreateRequestObject{
				Body: &api.CourseRegistrationsV1CreateJSONRequestBody{
					SubjectId: "s1",
				},
			})
			tt.validate(t, resp, err)
		})
	}
}

func TestCourseRegistrationsV1Delete(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		validate func(t *testing.T, resp api.CourseRegistrationsV1DeleteResponseObject, err error)
	}{
		{
			name:    "正常に履修登録を削除できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1DeleteResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.CourseRegistrationsV1Delete204Response)
				require.True(t, ok, "レスポンスが204レスポンスではありません")
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1DeleteResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1DeleteResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.CourseRegistrationsV1Delete401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("delete", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.CourseRegistrationsV1DeleteResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete course registration")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.CourseRegistrationsV1Delete(tt.ctx, api.CourseRegistrationsV1DeleteRequestObject{
				Id: "cr1",
			})
			tt.validate(t, resp, err)
		})
	}
}
