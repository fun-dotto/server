// Package main は Cloud Run Job として動作するマイグレーション実行エントリ。
//
// 役割:
//   - cloudsqlconn を使い IAM 認証付きで Cloud SQL (PostgreSQL) に接続する。
//   - ariga.io/atlas/sql/migrate (Executor) で ./migrations 配下の versioned SQL を
//     `atlas migrate apply` 相当として適用する。Cloud SQL Auth Proxy のサイドカーは
//     使わず、atlas CLI バイナリを同梱せずに *sql.DB 経由で完結させる。
//   - 外側から context.WithTimeout(10 分) を必ず掛け、ハング時は確実に non-zero exit させる。
//
// 環境変数:
//   - INSTANCE_CONNECTION_NAME: Cloud SQL の "<project>:<region>:<instance>"
//   - DB_NAME:                  接続先データベース名
//   - DB_IAM_USER:              Cloud SQL IAM ユーザ (例: "dotto-service@<project>.iam")
//   - MIGRATIONS_DIR (optional): デフォルト "/migrations" (コンテナ同梱パス)
//   - MIGRATE_TIMEOUT (optional): デフォルト "10m"
//
// NOTE: Phase 1 ではエントリポイントの骨組みのみ。Atlas のリビジョン管理テーブル
//       (atlas_schema_revisions) を扱う RevisionReadWriter の選定や、エラー時の
//       ロールバック方針は §E で実装を仕上げる。
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
)

func main() {
	timeout := envOrDefault("MIGRATE_TIMEOUT", "10m")
	d, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("MIGRATE_TIMEOUT のパースに失敗: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("migrate-job 失敗: %v", err)
	}
	log.Println("migrate-job 完了")
}

func run(ctx context.Context) error {
	instance := mustEnv("INSTANCE_CONNECTION_NAME")
	dbName := mustEnv("DB_NAME")
	iamUser := mustEnv("DB_IAM_USER")
	migrationsDir := envOrDefault("MIGRATIONS_DIR", "/migrations")

	db, cleanup, err := openCloudSQL(ctx, instance, iamUser, dbName)
	if err != nil {
		return err
	}
	defer cleanup()

	drv, err := postgres.Open(db)
	if err != nil {
		return fmt.Errorf("atlas postgres ドライバ生成失敗: %w", err)
	}

	dir, err := migrate.NewLocalDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrations ディレクトリのオープン失敗 (%s): %w", migrationsDir, err)
	}

	// TODO(§E): atlas_schema_revisions を使う RevisionReadWriter に差し替える。
	exec, err := migrate.NewExecutor(drv, dir, migrate.NopRevisionReadWriter{})
	if err != nil {
		return fmt.Errorf("migrate.Executor 生成失敗: %w", err)
	}

	if err := exec.ExecuteN(ctx, 0); err != nil {
		if err == migrate.ErrNoPendingFiles {
			log.Println("適用すべき差分なし")
			return nil
		}
		return fmt.Errorf("マイグレーション適用失敗: %w", err)
	}
	return nil
}

// openCloudSQL は cloudsqlconn ドライバを登録し、IAM 認証で Cloud SQL に接続する。
// 返り値の cleanup は defer で必ず呼ぶこと。
func openCloudSQL(ctx context.Context, instance, iamUser, dbName string) (*sql.DB, func(), error) {
	clean, err := pgxv5.RegisterDriver(
		"cloudsql-pgx",
		cloudsqlconn.WithIAMAuthN(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cloudsqlconn ドライバ登録失敗: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", instance, iamUser, dbName)
	db, err := sql.Open("cloudsql-pgx", dsn)
	if err != nil {
		_ = clean()
		return nil, nil, fmt.Errorf("sql.Open 失敗: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		_ = clean()
		return nil, nil, fmt.Errorf("DB Ping 失敗: %w", err)
	}

	return db, func() {
		_ = db.Close()
		if cerr := clean(); cerr != nil {
			log.Printf("cloudsqlconn cleanup 失敗: %v", cerr)
		}
	}, nil
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("環境変数 %s が未設定", k)
	}
	return v
}

func envOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
