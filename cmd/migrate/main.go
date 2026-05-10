// Package main は Cloud Run Job として起動し、cloudsqlconn 経由で
// IAM 認証された接続から `migrations/` 配下の SQL を適用するエントリポイント。
//
// diff の生成はローカル / CI で `atlas migrate diff --env local` が担当し、
// このバイナリは適用 (apply) のみを行う。
package main

import (
	"context"
	"log"

	"github.com/fun-dotto/server/internal/migrate"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	gormdb, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	sqlDB, err := gormdb.DB()
	if err != nil {
		log.Fatalf("failed to obtain *sql.DB: %v", err)
	}
	defer sqlDB.Close()

	const migrationsPath = "migrations"
	if err := migrate.ApplyMigrations(ctx, sqlDB, migrationsPath); err != nil {
		log.Fatalf("apply failed: %v", err)
	}
	log.Println("Migrations applied successfully.")
}
