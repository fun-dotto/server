// Command build-subjects-job は syllabuses から科目（subjects と子テーブル）を生成する。
// 科目が存在しない授業のみ作成し、既存の科目は更新しない。
//
// 環境変数:
//   - TARGET_YEAR: 対象年度（省略時は当年度）
package main

import (
	"context"
	"log"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/config"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/joho/godotenv"
)

func main() {
	logging.Setup()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	year, err := config.TargetYear()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	summary, err := service.NewSubjectBuildService(conn).Build(context.Background(), year)
	log.Printf("build summary: year=%d created=%d skipped=%d failed=%d", year, summary.Created, summary.Skipped, summary.Failed)
	if err != nil {
		log.Fatalf("Failed to build subjects: %v", err)
	}
}
