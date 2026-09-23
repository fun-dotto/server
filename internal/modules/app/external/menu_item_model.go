package external

import (
	funch_api "github.com/fun-dotto/server/gen/funch"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainMenuItem は外部APIのMenuItemをDomainに変換する
func ToDomainMenuItem(m funch_api.MenuItem) domain.MenuItem {
	prices := make([]domain.MenuPrice, len(m.Prices))
	for i, p := range m.Prices {
		prices[i] = domain.MenuPrice{
			Size:  domain.MenuSize(p.Size),
			Price: p.Price,
		}
	}
	return domain.MenuItem{
		ID:       m.Id,
		Date:     m.Date.Time,
		Name:     m.Name,
		ImageURL: m.ImageUrl,
		Category: domain.MenuCategory(m.Category),
		Prices:   prices,
	}
}

// ToExternalMenuItemQuery はDomainのMenuItemQueryを外部APIのパラメータに変換する
func ToExternalMenuItemQuery(q domain.MenuItemQuery) *funch_api.MenuItemsV1ListParams {
	return &funch_api.MenuItemsV1ListParams{
		Date: openapi_types.Date{Time: q.Date},
	}
}
