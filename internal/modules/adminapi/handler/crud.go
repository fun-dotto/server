package handler

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// crud は標準メソッドに共通する処理（永続化・ページネーション・エラー変換）をまとめる。
type crud[T any] struct {
	name resourceName
	db   *gorm.DB
	// keys はリソースの識別子を返す。順序がページネーションのソートキーになる。
	keys     func(T) []keyField
	preloads []string
	// createOmit は作成時に保存しない関連。nil ならすべての関連を保存しない。
	createOmit []string
	// updateFn は既定の更新処理を置き換える。子テーブルを持つリソースで使う。
	updateFn func(db *gorm.DB, key map[string]any, v *T, columns []string) error
}

func (c crud[T]) query(ctx context.Context) *gorm.DB {
	q := c.db.WithContext(ctx)
	for _, p := range c.preloads {
		q = q.Preload(p)
	}
	return q
}

func (c crud[T]) list(ctx context.Context, pageSize int32, pageToken string) ([]T, string, error) {
	size, err := normalizePageSize(pageSize)
	if err != nil {
		return nil, "", err
	}
	var zero T
	keys := c.keys(zero)
	after, err := decodePageToken(pageToken, len(keys))
	if err != nil {
		return nil, "", err
	}
	columns := make([]string, len(keys))
	for i, k := range keys {
		columns[i] = k.name
	}
	q := c.query(ctx)
	if after != nil {
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(after)), ",")
		q = q.Where(fmt.Sprintf("(%s) > (%s)", strings.Join(columns, ","), placeholders), after...)
	}
	for _, column := range columns {
		q = q.Order(column)
	}
	var items []T
	if err := q.Limit(size + 1).Find(&items).Error; err != nil {
		return nil, "", toConnectError(ctx, c.name, nil, err)
	}
	if len(items) <= size {
		return items, "", nil
	}
	items = items[:size]
	next, err := encodePageToken(cursorOf(c.keys(items[size-1])))
	if err != nil {
		return nil, "", toConnectError(ctx, c.name, nil, err)
	}
	return items, next, nil
}

func (c crud[T]) get(ctx context.Context, keys []keyField) (T, error) {
	var v T
	if err := c.query(ctx).Where(keyOf(keys)).Take(&v).Error; err != nil {
		return v, toConnectError(ctx, c.name, keys, err)
	}
	return v, nil
}

func (c crud[T]) create(ctx context.Context, v *T) (T, error) {
	omit := c.createOmit
	if omit == nil {
		omit = []string{clause.Associations}
	}
	if err := c.db.WithContext(ctx).Omit(omit...).Create(v).Error; err != nil {
		return *v, toConnectError(ctx, c.name, c.keys(*v), err)
	}
	// DB が補完した値（既定値・関連）を含めて返すため読み直す。
	return c.get(ctx, c.keys(*v))
}

// update は keys で特定したレコードの columns だけを v の値で更新する。
// columns には列名または構造体のフィールド名を指定する。
func (c crud[T]) update(ctx context.Context, keys []keyField, v *T, columns []string) (T, error) {
	if len(columns) > 0 {
		fn := c.updateFn
		if fn == nil {
			fn = updateColumns[T]
		}
		if err := fn(c.db.WithContext(ctx), keyOf(keys), v, columns); err != nil {
			return *v, toConnectError(ctx, c.name, keys, err)
		}
	}
	return c.get(ctx, keys)
}

func (c crud[T]) delete(ctx context.Context, keys []keyField) error {
	res := c.db.WithContext(ctx).Where(keyOf(keys)).Delete(new(T))
	if res.Error == nil && res.RowsAffected == 0 {
		res.Error = gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return toConnectError(ctx, c.name, keys, res.Error)
	}
	return nil
}

func updateColumns[T any](db *gorm.DB, key map[string]any, v *T, columns []string) error {
	res := db.Model(new(T)).Where(key).Select(columns).Omit(clause.Associations).Updates(v)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// cursorOf はページトークンに格納する値を JSON で表せる型（string / int64）にそろえる。
func cursorOf(keys []keyField) []any {
	out := make([]any, len(keys))
	for i, k := range keys {
		switch v := k.value.(type) {
		case int:
			out[i] = int64(v)
		case int32:
			out[i] = int64(v)
		case int64:
			out[i] = v
		default:
			out[i] = k.String()
		}
	}
	return out
}
