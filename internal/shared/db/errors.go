package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgresErrorCode は err に含まれる PostgreSQL のエラーコード（SQLSTATE）を返す。
// PostgreSQL のエラーでない場合は空文字列を返す。戻り値は github.com/jackc/pgerrcode の定数と比べる。
func PostgresErrorCode(err error) string {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code
	}
	return ""
}
