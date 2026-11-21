package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func TestPostgresConstraintName(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "ラップされた PgError", err: fmt.Errorf("create course registration: %w", &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: "fk_course_registrations_subject"}), want: "fk_course_registrations_subject"},
		{name: "PostgreSQL のエラーでなければ空文字列", err: errors.New("boom"), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.PostgresConstraintName(tt.err); got != tt.want {
				t.Errorf("PostgresConstraintName() = %q, want %q", got, tt.want)
			}
		})
	}
}
