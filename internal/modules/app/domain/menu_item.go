package domain

import "time"

// MenuCategory メニュー分類
type MenuCategory string

const (
	MenuCategorySetAndSingle MenuCategory = "SetAndSingle"
	MenuCategoryBowlAndCurry MenuCategory = "BowlAndCurry"
	MenuCategoryNoodle       MenuCategory = "Noodle"
	MenuCategorySide         MenuCategory = "Side"
	MenuCategoryDessert      MenuCategory = "Dessert"
)

// MenuSize メニューサイズ
type MenuSize string

const (
	MenuSizeSmall  MenuSize = "Small"
	MenuSizeMedium MenuSize = "Medium"
	MenuSizeLarge  MenuSize = "Large"
)

// MenuPrice メニュー価格
type MenuPrice struct {
	Size  MenuSize
	Price int32
}

// MenuItem メニュー項目
type MenuItem struct {
	ID       string
	Date     time.Time
	Name     string
	ImageURL string
	Category MenuCategory
	Prices   []MenuPrice
}

// MenuItemQuery メニュー検索クエリ
type MenuItemQuery struct {
	Date time.Time
}
