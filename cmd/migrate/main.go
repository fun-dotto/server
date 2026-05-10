package main

import (
	"context"
	"flag"
	"log"

	"github.com/fun-dotto/server/internal/migrate"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	mode := flag.String("mode", "all", "Execution mode: diff, apply, or all")
	flag.Parse()

	ctx := context.Background()

	// 1. Cloud SQL IAM接続のセットアップ
	gormdb, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	db, err := gormdb.DB()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	// 2. モードに応じた実行
	migrationsPath := "migrations"

	switch *mode {
	case "diff":
		if err := migrate.GenerateDiff(ctx, db, migrationsPath, "schema_update"); err != nil {
			log.Fatalf("diff failed: %v", err)
		}
		log.Println("Diff generated successfully.")

	case "apply":
		if err := migrate.ApplyMigrations(ctx, db, migrationsPath); err != nil {
			log.Fatalf("apply failed: %v", err)
		}
		log.Println("Migrations applied successfully.")

	case "all":
		// 差分を作ってから即適用するパターン
		if err := migrate.GenerateDiff(ctx, db, migrationsPath, "auto_gen"); err != nil {
			log.Fatalf("diff failed: %v", err)
		}
		if err := migrate.ApplyMigrations(ctx, db, migrationsPath); err != nil {
			log.Fatalf("apply failed: %v", err)
		}
		log.Println("All processes finished successfully.")
	}
}
