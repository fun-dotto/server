package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *CancelledClassRepository) ListByDate(ctx context.Context, date time.Time) ([]domain.CancelledClass, error) {
	var rows []model.CancelledClass
	if err := r.db.WithContext(ctx).
		Preload("Subject").
		Where("date = ?", date.Format(dateLayout)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CancelledClass, 0, len(rows))
	for i := range rows {
		out = append(out, cancelledClassToDomain(&rows[i]))
	}
	return out, nil
}
