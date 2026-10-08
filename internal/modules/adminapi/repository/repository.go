// Package repository は admin-api のデータアクセス層。GORM モデルをそのまま扱う。
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrAlreadyExists       = errors.New("record already exists")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrCheckViolation      = errors.New("check constraint violation")
)

// Key はレコードを特定する主キーの列名と値の組。
type Key map[string]any

// Repository は単一テーブルに対する CRUD を提供する。
type Repository[T any] struct {
	db *gorm.DB
	// keyColumns はページネーションのソートキー兼主キーの列名。
	keyColumns []string
	preloads   []string
}

func newRepository[T any](db *gorm.DB, keyColumns []string, preloads ...string) *Repository[T] {
	return &Repository[T]{db: db, keyColumns: keyColumns, preloads: preloads}
}

func (r *Repository[T]) query(ctx context.Context) *gorm.DB {
	q := r.db.WithContext(ctx)
	for _, p := range r.preloads {
		q = q.Preload(p)
	}
	return q
}

// List は keyColumns の昇順で、after より後ろのレコードを最大 limit 件返す。
// after は keyColumns と同じ順序の値で、nil なら先頭から返す。
func (r *Repository[T]) List(ctx context.Context, after []any, limit int) ([]T, error) {
	q := r.query(ctx)
	if after != nil {
		if len(after) != len(r.keyColumns) {
			return nil, fmt.Errorf("cursor has %d values, want %d", len(after), len(r.keyColumns))
		}
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(after)), ",")
		q = q.Where(fmt.Sprintf("(%s) > (%s)", strings.Join(r.keyColumns, ","), placeholders), after...)
	}
	for _, c := range r.keyColumns {
		q = q.Order(c)
	}
	var out []T
	if err := q.Limit(limit).Find(&out).Error; err != nil {
		return nil, translateError(err)
	}
	return out, nil
}

func (r *Repository[T]) Get(ctx context.Context, key Key) (T, error) {
	var out T
	if err := r.query(ctx).Where(map[string]any(key)).Take(&out).Error; err != nil {
		return out, translateError(err)
	}
	return out, nil
}

func (r *Repository[T]) Create(ctx context.Context, v *T) error {
	return translateError(r.db.WithContext(ctx).Omit(clause.Associations).Create(v).Error)
}

// Update は key で特定したレコードの columns だけを v の値で更新する。
// columns には列名または構造体のフィールド名を指定する。
func (r *Repository[T]) Update(ctx context.Context, key Key, v *T, columns []string) error {
	if len(columns) == 0 {
		_, err := r.Get(ctx, key)
		return err
	}
	return updateColumns[T](r.db.WithContext(ctx), key, v, columns)
}

func (r *Repository[T]) Delete(ctx context.Context, key Key) error {
	res := r.db.WithContext(ctx).Where(map[string]any(key)).Delete(new(T))
	if res.Error != nil {
		return translateError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func updateColumns[T any](db *gorm.DB, key Key, v *T, columns []string) error {
	res := db.Model(new(T)).Where(map[string]any(key)).Select(columns).Omit(clause.Associations).Updates(v)
	if res.Error != nil {
		return translateError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s", ErrAlreadyExists, pgErr.ConstraintName)
		case "23503":
			return fmt.Errorf("%w: %s", ErrForeignKeyViolation, pgErr.ConstraintName)
		case "23514":
			return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.ConstraintName)
		case "22P02", "22007", "22008":
			// 不正な UUID・日付など、DB が値を解釈できなかった場合
			return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Message)
		}
	}
	return err
}
