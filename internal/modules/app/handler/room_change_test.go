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

func TestRoomChangesV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		validate func(t *testing.T, resp api.RoomChangesV1ListResponseObject, err error)
	}{
		{
			name:    "正常に教室変更一覧が取得できる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.RoomChangesV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.RoomChangesV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.Len(t, result.RoomChanges, 1)
				assert.Equal(t, "room-change-1", result.RoomChanges[0].Id)
				assert.Equal(t, api.DottoFoundationV1Period("Period3"), result.RoomChanges[0].Period)
				assert.Equal(t, "s1", result.RoomChanges[0].Subject.Id)
				assert.Equal(t, "r1", result.RoomChanges[0].OriginalRoom.Id)
				assert.Equal(t, "101講義室", result.RoomChanges[0].OriginalRoom.Name)
				assert.Equal(t, api.DottoFoundationV1Floor("Floor1"), result.RoomChanges[0].OriginalRoom.Floor)
				assert.Equal(t, "r2", result.RoomChanges[0].NewRoom.Id)
				assert.Equal(t, "201講義室", result.RoomChanges[0].NewRoom.Name)
				assert.Equal(t, api.DottoFoundationV1Floor("Floor2"), result.RoomChanges[0].NewRoom.Floor)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.RoomChangesV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getRoomChanges", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.RoomChangesV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get room changes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.RoomChangesV1List(context.Background(), api.RoomChangesV1ListRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}
