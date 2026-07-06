package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// InsertBatch は faculty_rooms への一括 INSERT を 1 トランザクションで行う。
// (faculty_id, room_id, year) の UNIQUE 制約に違反した場合はロールバックしエラーを返す。
func (r *FacultyRoomRepository) InsertBatch(ctx context.Context, rows []domain.FacultyRoomInsert) error {
	if len(rows) == 0 {
		return nil
	}

	records := make([]model.FacultyRoom, len(rows))
	for i, row := range rows {
		records[i] = facultyRoomFromDomain(row.FacultyID, row.RoomID, row.Year)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&records).Error
	})
}
