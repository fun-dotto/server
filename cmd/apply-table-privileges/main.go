// Package main は Cloud Run Job として起動し、cloudsqlconn 経由で
// IAM 認証された接続から current_schema() 配下の全テーブルに対して
// 所有者の付け替えとロールごとの権限付与を再適用するエントリポイント。
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

	if err := migrate.ApplyTablePrivileges(ctx, sqlDB); err != nil {
		log.Fatalf("apply table privileges failed: %v", err)
	}
	log.Println("Table privileges applied successfully.")
}
