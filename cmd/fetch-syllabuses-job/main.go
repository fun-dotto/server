package main

import (
	"context"
	"flag"
	"log"
	"time"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/syllabusfetch"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/joho/godotenv"
)

func main() {
	logging.Setup()

	year := flag.Int("year", academicdomain.CurrentAcademicYear(time.Now()), "Academic year of syllabuses to fetch (default: current academic year in JST)")
	interval := flag.Duration("interval", time.Second, "Minimum interval between requests to the academic affairs system")
	flag.Parse()

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

	job := syllabusfetch.NewJob(syllabusfetch.NewClient(syllabusfetch.DefaultBaseURL, *interval), syllabusfetch.NewStore(conn))

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
