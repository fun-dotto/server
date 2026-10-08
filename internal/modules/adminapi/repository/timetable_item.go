package repository

import (
	"context"
	"slices"

	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// TimetableItemRepository は教室（timetable_item_rooms）を含めて TimetableItem を扱う。
type TimetableItemRepository struct {
	*Repository[model.TimetableItem]
}

func NewTimetableItemRepository(db *gorm.DB) *TimetableItemRepository {
	return &TimetableItemRepository{
		Repository: newRepository[model.TimetableItem](db, []string{"id"}, "Rooms"),
	}
}

func (r *TimetableItemRepository) Create(ctx context.Context, v *model.TimetableItem) error {
	return translateError(r.db.WithContext(ctx).Omit("Subject", "Rooms.Room").Create(v).Error)
}

// Update は columns に Rooms が含まれる場合、教室の行をすべて v の内容で置き換える。
func (r *TimetableItemRepository) Update(ctx context.Context, key Key, v *model.TimetableItem, columns []string) error {
	if len(columns) == 0 {
		_, err := r.Get(ctx, key)
		return err
	}
	return translateError(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cols := slices.DeleteFunc(slices.Clone(columns), func(c string) bool { return c == "Rooms" })
		cols = append(cols, "UpdatedAt")
		if err := updateColumns[model.TimetableItem](tx, key, v, cols); err != nil {
			return err
		}
		if !slices.Contains(columns, "Rooms") {
			return nil
		}
		var current model.TimetableItem
		if err := tx.Where(map[string]any(key)).Take(&current).Error; err != nil {
			return err
		}
		if err := tx.Where("timetable_item_id = ?", current.ID.String()).Delete(&model.TimetableItemRoom{}).Error; err != nil {
			return err
		}
		if len(v.Rooms) == 0 {
			return nil
		}
		rooms := slices.Clone(v.Rooms)
		for i := range rooms {
			rooms[i].TimetableItemID = current.ID.String()
		}
		return tx.Omit("Room").Create(&rooms).Error
	}))
}
