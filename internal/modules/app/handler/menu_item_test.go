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

func TestMenuItemsV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		validate func(t *testing.T, resp api.MenuItemsV1ListResponseObject, err error)
	}{
		{
			name:    "正常にメニュー一覧が取得できる",
			handler: NewHandler(WithFunchService(service.NewFunchService(repository.NewMockFunchRepository()))),
			validate: func(t *testing.T, resp api.MenuItemsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.MenuItemsV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				require.Len(t, result.MenuItems, 1)
				assert.Equal(t, "menu-1", result.MenuItems[0].Id)
				assert.Equal(t, "からあげ定食", result.MenuItems[0].Name)
				assert.Equal(t, api.Category("SetAndSingle"), result.MenuItems[0].Category)
				require.Len(t, result.MenuItems[0].Prices, 1)
				assert.Equal(t, api.Size("Medium"), result.MenuItems[0].Prices[0].Size)
				assert.Equal(t, int32(500), result.MenuItems[0].Prices[0].Price)
			},
		},
		{
			name:    "funchServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.MenuItemsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errFunchServiceNotConfigured)
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithFunchService(service.NewFunchService(repository.NewMockFunchRepositoryWithError("getMenuItems", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.MenuItemsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get menu items")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.MenuItemsV1List(context.Background(), api.MenuItemsV1ListRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}
