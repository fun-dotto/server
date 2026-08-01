package repository

import (
	"context"
	"fmt"

	funch_api "github.com/fun-dotto/server/gen/funch"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/external"
)

type FunchRepository struct {
	client *funch_api.ClientWithResponses
}

func NewFunchRepository(client *funch_api.ClientWithResponses) *FunchRepository {
	return &FunchRepository{client: client}
}

// GetMenuItems は外部APIからメニュー一覧を取得する
func (r *FunchRepository) GetMenuItems(query domain.MenuItemQuery) ([]domain.MenuItem, error) {
	params := external.ToExternalMenuItemQuery(query)

	response, err := r.client.MenuItemsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to call funch API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get menu items: status %d", response.StatusCode())
	}

	items := response.JSON200.MenuItems
	result := make([]domain.MenuItem, len(items))
	for i, item := range items {
		result[i] = external.ToDomainMenuItem(item)
	}

	return result, nil
}
