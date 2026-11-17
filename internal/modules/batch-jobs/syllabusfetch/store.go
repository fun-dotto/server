package syllabusfetch

import (
	"context"
	"errors"
	"time"

	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// SaveResult は Store.Save で行った操作。
type SaveResult int

const (
	SaveInserted SaveResult = iota
	SaveUpdated
	SaveUnchanged
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Save は同じ (lesson_id, year) の行に取得結果を反映する。
// シラバスと科目を 1:1 に保つため新しいバージョン行は作らず、既存行の id を保ったまま上書きする。
// 内容が同じ場合は最終確認日時として fetched_at だけを更新する。
func (s *Store) Save(ctx context.Context, syllabus model.Syllabus, fetchedAt time.Time) (SaveResult, error) {
	syllabus.ContentHash = syllabus.ComputeContentHash()
	syllabus.FetchedAt = fetchedAt

	var result SaveResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Syllabus
		err := tx.Where("lesson_id = ? AND year = ?", syllabus.LessonID, syllabus.Year).
			Order("fetched_at DESC").
			First(&existing).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			result = SaveInserted
			return tx.Create(&syllabus).Error
		case err != nil:
			return err
		case existing.ContentHash == syllabus.ContentHash:
			result = SaveUnchanged
			return tx.Model(&existing).Update("fetched_at", fetchedAt).Error
		default:
			result = SaveUpdated
			syllabus.ID = existing.ID
			syllabus.CreatedAt = existing.CreatedAt
			// Select("*") で空文字になった列も上書きする（Updates は既定でゼロ値を無視する）
			return tx.Model(&existing).Select("*").Omit("id", "created_at").Updates(&syllabus).Error
		}
	})
	return result, err
}
