package handler

import (
	"context"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimetableItemsV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		request  api.TimetableItemsV1ListRequestObject
		validate func(t *testing.T, resp api.TimetableItemsV1ListResponseObject, err error)
	}{
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			request: api.TimetableItemsV1ListRequestObject{},
			validate: func(t *testing.T, resp api.TimetableItemsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "academicService is not configured")
			},
		},
		{
			name: "正常に時間割アイテム一覧を取得できる",
			handler: NewHandler(
				WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository())),
			),
			request: api.TimetableItemsV1ListRequestObject{
				Params: api.TimetableItemsV1ListParams{
					Semesters: []api.DottoFoundationV1CourseSemester{api.Q1},
				},
			},
			validate: func(t *testing.T, resp api.TimetableItemsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.TimetableItemsV1List200JSONResponse)
				require.True(t, ok)
				assert.Len(t, result.TimetableItems, 1)
				assert.Equal(t, "t1", result.TimetableItems[0].Id)
				require.NotNil(t, result.TimetableItems[0].Slot)
				assert.Len(t, result.TimetableItems[0].Rooms, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.TimetableItemsV1List(context.Background(), tt.request)
			tt.validate(t, resp, err)
		})
	}
}
