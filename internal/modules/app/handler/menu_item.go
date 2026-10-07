package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/app"
)

func (h *Handler) MenuItemsV1List(_ context.Context, _ api.MenuItemsV1ListRequestObject) (api.MenuItemsV1ListResponseObject, error) {
	return nil, errors.New("menu items are not available")
}
