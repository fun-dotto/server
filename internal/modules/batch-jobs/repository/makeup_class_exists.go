package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/shared/model"
)

// Exists は (subject_id, date, period) の自然キーで既存レコードの有無を返す。
func (r *MakeupClassRepository) Exists(ctx context.Context, subjectID string, date time.Time, period string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.MakeupClass{}).
		Where("subject_id = ? AND date = ? AND period = ?", parseUUIDOrNil(subjectID), date.Format(dateLayout), period).
		Count(&count).Error
	return count > 0, err
}
