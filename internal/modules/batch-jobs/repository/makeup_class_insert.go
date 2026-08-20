package repository

import (
	"context"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
)

// Insert は補講レコードを 1 件挿入する。自然キーの重複チェックは呼び出し側（Exists）で行う。
func (r *MakeupClassRepository) Insert(ctx context.Context, mc academicdomain.MakeupClass) (academicdomain.MakeupClass, error) {
	record := makeupClassFromDomain(mc)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return academicdomain.MakeupClass{}, err
	}
	return makeupClassToDomain(record), nil
}
