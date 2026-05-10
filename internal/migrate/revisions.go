package migrate

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ariga.io/atlas/sql/migrate"
)

// pgRevisions は atlas CLI 互換の `atlas_schema_revisions` テーブルを
// 直接 SQL で読み書きする migrate.RevisionReadWriter 実装。
//
// テーブルスキーマは atlas CLI が ent で作成するものと一致させ、
// ローカルで `atlas migrate set` などの CLI 操作と相互運用できるようにする。
type pgRevisions struct {
	db     *sql.DB
	schema string
	table  string
}

const (
	revisionsSchema = "atlas_schema_revisions"
	revisionsTable  = "atlas_schema_revisions"
)

func newPGRevisions(ctx context.Context, db *sql.DB) (*pgRevisions, error) {
	r := &pgRevisions{db: db, schema: revisionsSchema, table: revisionsTable}
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *pgRevisions) ensureTable(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS "atlas_schema_revisions"`); err != nil {
		return err
	}
	const stmt = `
CREATE TABLE IF NOT EXISTS "atlas_schema_revisions"."atlas_schema_revisions" (
    version          VARCHAR PRIMARY KEY NOT NULL,
    description      VARCHAR NOT NULL,
    type             BIGINT NOT NULL DEFAULT 2,
    applied          BIGINT NOT NULL DEFAULT 0,
    total            BIGINT NOT NULL DEFAULT 0,
    executed_at      TIMESTAMPTZ NOT NULL,
    execution_time   BIGINT NOT NULL,
    error            TEXT,
    error_stmt       TEXT,
    hash             VARCHAR NOT NULL,
    partial_hashes   JSONB,
    operator_version VARCHAR NOT NULL
);`
	_, err := r.db.ExecContext(ctx, stmt)
	return err
}

func (r *pgRevisions) Ident() *migrate.TableIdent {
	return &migrate.TableIdent{Name: r.table, Schema: r.schema}
}

func (r *pgRevisions) qualified() string {
	return fmt.Sprintf(`"%s"."%s"`, r.schema, r.table)
}

func (r *pgRevisions) ReadRevisions(ctx context.Context) ([]*migrate.Revision, error) {
	q := `SELECT version, description, type, applied, total, executed_at,
		execution_time, error, error_stmt, hash, partial_hashes, operator_version
		FROM ` + r.qualified() + ` ORDER BY executed_at`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("read revisions: %w", err)
	}
	defer rows.Close()

	var revs []*migrate.Revision
	for rows.Next() {
		rev, err := scanRevision(rows)
		if err != nil {
			return nil, err
		}
		revs = append(revs, rev)
	}
	return revs, rows.Err()
}

func (r *pgRevisions) ReadRevision(ctx context.Context, version string) (*migrate.Revision, error) {
	q := `SELECT version, description, type, applied, total, executed_at,
		execution_time, error, error_stmt, hash, partial_hashes, operator_version
		FROM ` + r.qualified() + ` WHERE version = $1`
	rev, err := scanRevision(r.db.QueryRowContext(ctx, q, version))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, migrate.ErrRevisionNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("read revision %q: %w", version, err)
	}
	return rev, nil
}

func (r *pgRevisions) WriteRevision(ctx context.Context, rev *migrate.Revision) error {
	partial, err := marshalPartialHashes(rev.PartialHashes)
	if err != nil {
		return err
	}
	stmt := `
INSERT INTO ` + r.qualified() + ` (
    version, description, type, applied, total, executed_at,
    execution_time, error, error_stmt, hash, partial_hashes, operator_version
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (version) DO UPDATE SET
    description      = EXCLUDED.description,
    type             = EXCLUDED.type,
    applied          = EXCLUDED.applied,
    total            = EXCLUDED.total,
    executed_at      = EXCLUDED.executed_at,
    execution_time   = EXCLUDED.execution_time,
    error            = EXCLUDED.error,
    error_stmt       = EXCLUDED.error_stmt,
    hash             = EXCLUDED.hash,
    partial_hashes   = EXCLUDED.partial_hashes,
    operator_version = EXCLUDED.operator_version`
	_, err = r.db.ExecContext(ctx, stmt,
		rev.Version, rev.Description, int64(rev.Type), int64(rev.Applied), int64(rev.Total),
		rev.ExecutedAt.UTC(), int64(rev.ExecutionTime), nullableString(rev.Error),
		nullableString(rev.ErrorStmt), rev.Hash, partial, rev.OperatorVersion,
	)
	if err != nil {
		return fmt.Errorf("write revision %q: %w", rev.Version, err)
	}
	return nil
}

func (r *pgRevisions) DeleteRevision(ctx context.Context, version string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM `+r.qualified()+` WHERE version = $1`, version); err != nil {
		return fmt.Errorf("delete revision %q: %w", version, err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRevision(s rowScanner) (*migrate.Revision, error) {
	var (
		rev        migrate.Revision
		typ        int64
		applied    int64
		total      int64
		execTime   int64
		executedAt time.Time
		errStr     sql.NullString
		errStmt    sql.NullString
		partial    sql.NullString
	)
	if err := s.Scan(
		&rev.Version, &rev.Description, &typ, &applied, &total, &executedAt,
		&execTime, &errStr, &errStmt, &rev.Hash, &partial, &rev.OperatorVersion,
	); err != nil {
		return nil, err
	}
	rev.Type = migrate.RevisionType(typ)
	rev.Applied = int(applied)
	rev.Total = int(total)
	rev.ExecutedAt = executedAt
	rev.ExecutionTime = time.Duration(execTime)
	rev.Error = errStr.String
	rev.ErrorStmt = errStmt.String
	if partial.Valid && partial.String != "" {
		if err := json.Unmarshal([]byte(partial.String), &rev.PartialHashes); err != nil {
			return nil, fmt.Errorf("unmarshal partial_hashes for %q: %w", rev.Version, err)
		}
	}
	return &rev, nil
}

func marshalPartialHashes(hashes []string) (any, error) {
	if hashes == nil {
		return nil, nil
	}
	b, err := json.Marshal(hashes)
	if err != nil {
		return nil, fmt.Errorf("marshal partial_hashes: %w", err)
	}
	return string(b), nil
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
