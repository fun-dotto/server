// Package migrate は、atlas CLI が生成した versioned SQL を本番環境
// (Cloud Run Job) で適用するためのランタイム実装。
//
// diff 生成はローカル / CI で `atlas migrate diff --env local` が担当し、
// このパッケージは cloudsqlconn 経由で IAM 認証された *sql.DB を受け取って
// `migrations/` 配下を Atlas Go API で適用するだけに専念する。
package migrate

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
)

// ApplyMigrations は dirPath 配下の未適用 SQL マイグレーションを順に DB へ適用する。
// 適用済みリビジョンは atlas CLI 互換の atlas_schema_revisions テーブルで管理する。
func ApplyMigrations(ctx context.Context, db *sql.DB, dirPath string) error {
	drv, err := postgres.Open(db)
	if err != nil {
		return fmt.Errorf("open postgres driver: %w", err)
	}

	dir, err := migrate.NewLocalDir(dirPath)
	if err != nil {
		return fmt.Errorf("open migration dir %q: %w", dirPath, err)
	}

	rrw, err := newPGRevisions(ctx, db)
	if err != nil {
		return fmt.Errorf("init revisions table: %w", err)
	}

	exec, err := migrate.NewExecutor(drv, dir, rrw)
	if err != nil {
		return fmt.Errorf("create executor: %w", err)
	}

	if err := exec.ExecuteN(ctx, 0); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
