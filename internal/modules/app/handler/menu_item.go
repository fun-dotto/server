package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) MenuItemsV1List(_ context.Context, request api.MenuItemsV1ListRequestObject) (api.MenuItemsV1ListResponseObject, error) {
	if h.funchService == nil {
		return nil, errFunchServiceNotConfigured
	}

	query := domain.MenuItemQuery{Date: request.Params.Date.Time}

	items, err := h.funchService.GetMenuItems(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items: %w", err)
	}

	apiItems := make([]api.MenuItem, len(items))
	for i, item := range items {
		apiItems[i] = toApiMenuItem(item)
	}

	return api.MenuItemsV1List200JSONResponse{
		MenuItems: apiItems,
	}, nil
}

func toApiMenuItem(m domain.MenuItem) api.MenuItem {
	prices := make([]api.Price, len(m.Prices))
	for i, p := range m.Prices {
		prices[i] = api.Price{
			Size:  api.Size(p.Size),
			Price: p.Price,
		}
	}
	return api.MenuItem{
		Id:       m.ID,
		Date:     openapi_types.Date{Time: m.Date},
		Name:     m.Name,
		ImageUrl: m.ImageURL,
		Category: api.Category(m.Category),
		Prices:   prices,
	}
}
