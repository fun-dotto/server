package repository

import (
	"context"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
)

// Insert は休講レコードを 1 件挿入する。自然キーの重複チェックは呼び出し側（Exists）で行う。
func (r *CancelledClassRepository) Insert(ctx context.Context, cc academicdomain.CancelledClass) (academicdomain.CancelledClass, error) {
	record := cancelledClassFromDomain(cc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return academicdomain.CancelledClass{}, err
	}
	return cancelledClassToDomain(record), nil
}
