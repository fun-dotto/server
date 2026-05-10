package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *RoomChangeRepository) ListByDate(ctx context.Context, date time.Time) ([]domain.RoomChange, error) {
	var rows []model.RoomChange
	if err := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("NewRoom").
		Where("date = ?", date.Format(dateLayout)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RoomChange, 0, len(rows))
	for i := range rows {
		out = append(out, roomChangeToDomain(&rows[i]))
	}
	return out, nil
}
