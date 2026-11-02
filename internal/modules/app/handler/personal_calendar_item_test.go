package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonalCalendarItemsV1List(t *testing.T) {
	dates := []openapi_types.Date{{Time: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)}}

	tests := []struct {
		name     string
		ctx      context.Context
		handler  *Handler
		validate func(t *testing.T, resp api.PersonalCalendarItemsV1ListResponseObject, err error)
	}{
		{
			name:    "正常に個人カレンダーアイテム一覧を取得できる",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.PersonalCalendarItemsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.PersonalCalendarItemsV1List200JSONResponse)
				require.True(t, ok)
				assert.Len(t, result.PersonalCalendarItems, 1)
				assert.Equal(t, api.DottoFoundationV1Period("Period1"), result.PersonalCalendarItems[0].Period)
				assert.Equal(t, "s1", result.PersonalCalendarItems[0].Subject.Id)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.PersonalCalendarItemsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "コンテキストにユーザーIDがない場合401を返す",
			ctx:     context.Background(),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.PersonalCalendarItemsV1ListResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.PersonalCalendarItemsV1List401Response)
				require.True(t, ok, "レスポンスが401レスポンスではありません")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			ctx:     ctxWithUserID("user1"),
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getPersonalCalendarItems", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.PersonalCalendarItemsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get personal calendar items")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.PersonalCalendarItemsV1List(tt.ctx, api.PersonalCalendarItemsV1ListRequestObject{
				Params: api.PersonalCalendarItemsV1ListParams{
					Dates: dates,
				},
			})
			tt.validate(t, resp, err)
		})
	}
}
