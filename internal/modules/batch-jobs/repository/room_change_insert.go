package repository

import (
	"context"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
)

// Insert は部屋変更レコードを 1 件挿入する。自然キーの重複チェックは呼び出し側（Exists）で行う。
func (r *RoomChangeRepository) Insert(ctx context.Context, rc academicdomain.RoomChange) (academicdomain.RoomChange, error) {
	record := roomChangeFromDomain(rc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return academicdomain.RoomChange{}, err
	}
	return roomChangeToDomain(record), nil
}
