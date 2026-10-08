package handler

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/adminapi/repository"
)

// Repository は Handler が依存するリソースごとの永続化の操作。
type Repository[T any] interface {
	List(ctx context.Context, after []any, limit int) ([]T, error)
	Get(ctx context.Context, key repository.Key) (T, error)
	Create(ctx context.Context, v *T) error
	Update(ctx context.Context, key repository.Key, v *T, columns []string) error
	Delete(ctx context.Context, key repository.Key) error
}

// crud は標準メソッドに共通する処理（ページネーション・エラー変換）をまとめる。
type crud[T any] struct {
	name resourceName
	repo Repository[T]
	// keys はリソースの識別子を返す。順序はページネーションのソートキーと一致させる。
	keys func(T) []keyField
}

func (c crud[T]) list(ctx context.Context, pageSize int32, pageToken string) ([]T, string, error) {
	size, err := normalizePageSize(pageSize)
	if err != nil {
		return nil, "", err
	}
	var zero T
	after, err := decodePageToken(pageToken, len(c.keys(zero)))
	if err != nil {
		return nil, "", err
	}
	items, err := c.repo.List(ctx, after, size+1)
	if err != nil {
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
	v, err := c.repo.Get(ctx, keyOf(keys))
	if err != nil {
		return v, toConnectError(ctx, c.name, keys, err)
	}
	return v, nil
}

func (c crud[T]) create(ctx context.Context, v *T) (T, error) {
	if err := c.repo.Create(ctx, v); err != nil {
		return *v, toConnectError(ctx, c.name, c.keys(*v), err)
	}
	// DB が補完した値（既定値・関連）を含めて返すため読み直す。
	return c.get(ctx, c.keys(*v))
}

func (c crud[T]) update(ctx context.Context, keys []keyField, v *T, columns []string) (T, error) {
	if err := c.repo.Update(ctx, keyOf(keys), v, columns); err != nil {
		return *v, toConnectError(ctx, c.name, keys, err)
	}
	return c.get(ctx, keys)
}

func (c crud[T]) delete(ctx context.Context, keys []keyField) error {
	if err := c.repo.Delete(ctx, keyOf(keys)); err != nil {
		return toConnectError(ctx, c.name, keys, err)
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
