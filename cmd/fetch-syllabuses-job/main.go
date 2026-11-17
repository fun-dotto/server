package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/syllabusfetch"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/joho/godotenv"
)

func main() {
	year := flag.Int("year", currentAcademicYear(time.Now()), "Academic year of syllabuses to fetch (default: current academic year in JST)")
	interval := flag.Duration("interval", time.Second, "Minimum interval between requests to the academic affairs system")
	flag.Parse()

	logging.Setup()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
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

	job := syllabusfetch.NewJob(syllabusfetch.NewClient(*interval), syllabusfetch.NewStore(conn))

	summary, err := job.Run(context.Background(), *year)
	if err != nil {
		log.Fatalf("Failed to fetch syllabuses: %v", err)
	}
	log.Printf(
		"fetch summary: year=%d listed=%d inserted=%d updated=%d unchanged=%d failed=%d",
		*year,
		summary.Listed,
		summary.Inserted,
		summary.Updated,
		summary.Unchanged,
		summary.Failed,
	)
	if summary.Failed > 0 {
		log.Fatalf("Failed to fetch %d syllabuses", summary.Failed)
	}
}

// currentAcademicYear は日本の学年度（4 月始まり）を返す。1〜3 月は前年の年度になる。
func currentAcademicYear(now time.Time) int {
	jst := now.In(time.FixedZone("Asia/Tokyo", 9*60*60))
	if jst.Month() < time.April {
		return jst.Year() - 1
	}
	return jst.Year()
}
