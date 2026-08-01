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

func TestCancelledClassesV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		validate func(t *testing.T, resp api.CancelledClassesV1ListResponseObject, err error)
	}{
		{
			name:    "正常に休講一覧が取得できる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.CancelledClassesV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.CancelledClassesV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.Len(t, result.CancelledClasses, 1)
				assert.Equal(t, "cancelled-1", result.CancelledClasses[0].Id)
				assert.Equal(t, "教員出張のため", result.CancelledClasses[0].Comment)
				assert.Equal(t, api.DottoFoundationV1Period("Period1"), result.CancelledClasses[0].Period)
				assert.Equal(t, "s1", result.CancelledClasses[0].Subject.Id)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.CancelledClassesV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getCancelledClasses", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.CancelledClassesV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get cancelled classes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.CancelledClassesV1List(context.Background(), api.CancelledClassesV1ListRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}
