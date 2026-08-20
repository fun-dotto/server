package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/server/internal/shared/model"
)

// Exists は (subject_id, date, period, original_room_id, new_room_id) の自然キーで
// 既存レコードの有無を返す。
func (r *RoomChangeRepository) Exists(ctx context.Context, subjectID string, date time.Time, period, originalRoomID, newRoomID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RoomChange{}).
		Where(
			"subject_id = ? AND date = ? AND period = ? AND original_room_id = ? AND new_room_id = ?",
			parseUUIDOrNil(subjectID), date.Format(dateLayout), period, parseUUIDOrNil(originalRoomID), parseUUIDOrNil(newRoomID),
		).
		Count(&count).Error
	return count > 0, err
}
