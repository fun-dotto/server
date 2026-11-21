package main

import (
	"context"
	"flag"
	"log"
	"time"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/repository"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/joho/godotenv"
)

func main() {
	logging.Setup()

	year := flag.Int("year", academicdomain.CurrentAcademicYear(time.Now()), "Academic year of syllabuses to build subjects from (default: current academic year in JST)")
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

	svc := service.NewSubjectBuildService(repository.NewSubjectBuildRepository(conn))

	summary, err := svc.BuildSubjects(context.Background(), *year)
	if err != nil {
		log.Fatalf("Failed to build subjects: %v", err)
	}
	log.Printf(
		"build subjects summary: year=%d created=%d skipped=%d failed=%d",
		*year,
		summary.Created,
		summary.Skipped,
		summary.Failed,
	)
	if summary.Failed > 0 {
		log.Fatalf("Failed to create %d subjects", summary.Failed)
	}
}
