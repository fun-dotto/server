package repository

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// MockFunchRepository はテスト用のモック
type MockFunchRepository struct {
	menuItems       []domain.MenuItem
	getMenuItemsErr error
}

func NewMockFunchRepository() *MockFunchRepository {
	return &MockFunchRepository{
		menuItems: []domain.MenuItem{
			{
				ID:       "menu-1",
				Date:     time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
				Name:     "からあげ定食",
				ImageURL: "https://example.com/karaage.png",
				Category: domain.MenuCategorySetAndSingle,
				Prices: []domain.MenuPrice{
					{Size: domain.MenuSizeMedium, Price: 500},
				},
			},
		},
	}
}

func NewMockFunchRepositoryWithError(field string, err error) *MockFunchRepository {
	m := NewMockFunchRepository()
	switch field {
	case "getMenuItems":
		m.getMenuItemsErr = err
	}
	return m
}

func (m *MockFunchRepository) GetMenuItems(_ domain.MenuItemQuery) ([]domain.MenuItem, error) {
	if m.getMenuItemsErr != nil {
		return nil, m.getMenuItemsErr
	}
	return m.menuItems, nil
}
