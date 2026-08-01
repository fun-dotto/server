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

func TestReservationsV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		validate func(t *testing.T, resp api.ReservationsV1ListResponseObject, err error)
	}{
		{
			name:    "正常に予約一覧が取得できる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.ReservationsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.ReservationsV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.Len(t, result.Reservations, 1)
				assert.Equal(t, "reservation-1", result.Reservations[0].Id)
				assert.Equal(t, "ゼミ", result.Reservations[0].Title)
				assert.Equal(t, "r1", result.Reservations[0].Room.Id)
				assert.Equal(t, api.DottoFoundationV1Floor("Floor1"), result.Reservations[0].Room.Floor)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.ReservationsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getReservations", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.ReservationsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get reservations")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.ReservationsV1List(context.Background(), api.ReservationsV1ListRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}
