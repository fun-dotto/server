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

// PostgresConstraintName は err に含まれる PostgreSQL のエラーで違反した制約の名前を返す。
// PostgreSQL のエラーでない場合や、制約違反でない場合は空文字列を返す。
func PostgresConstraintName(err error) string {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.ConstraintName
	}
	return ""
}
