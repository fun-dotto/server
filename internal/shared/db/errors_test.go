package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/fun-dotto/server/internal/shared/db"
)

func TestPostgresErrorCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "PgError", err: &pgconn.PgError{Code: pgerrcode.UniqueViolation}, want: pgerrcode.UniqueViolation},
		{name: "ラップされた PgError", err: fmt.Errorf("create subject: %w", &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation}), want: pgerrcode.ForeignKeyViolation},
		{name: "PostgreSQL のエラーでなければ空文字列", err: errors.New("boom"), want: ""},
		{name: "nil は空文字列", err: nil, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.PostgresErrorCode(tt.err); got != tt.want {
				t.Errorf("PostgresErrorCode() = %q, want %q", got, tt.want)
			}
		})
	}
}
