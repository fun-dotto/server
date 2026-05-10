package migrate

import (
	"context"
	"database/sql"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
)

// ApplyMigrations はディレクトリ内のSQLファイルをDBに適用します
func ApplyMigrations(ctx context.Context, db *sql.DB, dirPath string) error {
	// 1. Atlasドライバの準備
	drv, err := postgres.Open(db)
	if err != nil {
		return err
	}

	// 2. マイグレーションファイルの読み込み
	dir, err := migrate.NewLocalDir(dirPath)
	if err != nil {
		return err
	}

	// 3. リビジョン管理テーブル (atlas_schema_revisions) の準備
	// これにより、適用済みのファイルが二重に実行されるのを防ぎます
	reg, err := migrate.NewEntRevisions(ctx, drv)
	if err != nil {
		return err
	}

	// 4. 実行器の作成と適用
	executor, err := migrate.NewExecutor(drv, dir, reg)
	if err != nil {
		return err
	}

	// 未適用のファイルをすべて実行
	return executor.ApplyAll(ctx)
}
