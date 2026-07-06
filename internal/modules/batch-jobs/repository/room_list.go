package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListRooms は教室名照合に使う (id, name) を全件返す。
func (r *RoomRepository) ListRooms(ctx context.Context) ([]domain.RoomRef, error) {
	var records []model.Room
	if err := r.db.WithContext(ctx).Select("id", "name").Find(&records).Error; err != nil {
		return nil, err
	}

	refs := make([]domain.RoomRef, len(records))
	for i, rec := range records {
		refs[i] = domain.RoomRef{ID: rec.ID.String(), Name: rec.Name}
	}
	return refs, nil
}
