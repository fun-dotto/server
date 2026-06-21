package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *MakeupClassRepository) ListByDate(ctx context.Context, date time.Time) ([]domain.MakeupClass, error) {
	var rows []model.MakeupClass
	if err := r.db.WithContext(ctx).
		Preload("Subject").
		Where("date = ?", date.Format(dateLayout)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.MakeupClass, 0, len(rows))
	for i := range rows {
		makeupClass, err := makeupClassToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, makeupClass)
	}
	return out, nil
}
