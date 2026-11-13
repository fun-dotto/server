// Command fetch-syllabuses-job は教務システムからシラバスを取得し、syllabuses に保存する。
//
// 環境変数:
//   - TARGET_YEAR: 対象年度（省略時は当年度）
//   - SYLLABUS_FETCH_INTERVAL: 教務システムへのリクエスト間隔（省略時は 1s）
//   - SLACK_WEBHOOK_URL: 差分を通知する Slack Incoming Webhook の URL（省略時は通知しない）
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/config"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/fun-dotto/server/internal/shared/slack"
	"github.com/fun-dotto/server/internal/shared/syllabussource"
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
	interval := time.Second
	if v := os.Getenv("SYLLABUS_FETCH_INTERVAL"); v != "" {
		if interval, err = time.ParseDuration(v); err != nil {
			log.Fatalf("invalid SYLLABUS_FETCH_INTERVAL: %v", err)
		}
	}
	var notifier service.Notifier
	if url := os.Getenv("SLACK_WEBHOOK_URL"); url != "" {
		notifier = slack.NewWebhook(url)
	} else {
		slog.Warn("SLACK_WEBHOOK_URL is not set; syllabus changes will not be notified")
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

	svc := service.NewSyllabusFetchService(conn, syllabussource.NewClient(interval), notifier)
	summary, err := svc.Fetch(context.Background(), year)
	log.Printf(
		"fetch summary: year=%d listed=%d inserted=%d unchanged=%d failed=%d changes=%d",
		year, summary.Listed, summary.Inserted, summary.Unchanged, summary.Failed, len(summary.Changes),
	)
	if err != nil {
		log.Fatalf("Failed to fetch syllabuses: %v", err)
	}
}
