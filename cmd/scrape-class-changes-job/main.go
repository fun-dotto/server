package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/repository"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/scraper"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	userID := os.Getenv("SCRAPE_USER_ID")
	password := os.Getenv("SCRAPE_USER_PASSWORD")
	if userID == "" || password == "" {
		log.Fatal("SCRAPE_USER_ID and SCRAPE_USER_PASSWORD must be set")
	}

	year := scraper.CurrentNendo(time.Now())
	if v := os.Getenv("SCRAPE_YEAR"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid SCRAPE_YEAR %q: %v", v, err)
		}
		year = parsed
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

	client, err := scraper.NewClient(year)
	if err != nil {
		log.Fatalf("Failed to create scraper client: %v", err)
	}

	svc := service.NewClassChangeScrapeService(
		client,
		repository.NewSubjectRepository(conn),
		repository.NewRoomRepository(conn),
		repository.NewCancelledClassRepository(conn),
		repository.NewMakeupClassRepository(conn),
		repository.NewRoomChangeRepository(conn),
	)

	summary, err := svc.Run(context.Background(), userID, password, year)
	if err != nil {
		log.Fatalf("Failed to scrape class changes: %v", err)
	}
	log.Printf(
		"scrape summary: year=%d cancelled(inserted=%d duplicated=%d skipped=%d) makeup(inserted=%d duplicated=%d skipped=%d) room_change(inserted=%d duplicated=%d skipped=%d) unmatched(subjects=%d rooms=%d)",
		year,
		summary.CancelledInserted, summary.CancelledDuplicated, summary.CancelledSkipped,
		summary.MakeupInserted, summary.MakeupDuplicated, summary.MakeupSkipped,
		summary.RoomChangeInserted, summary.RoomChangeDuplicated, summary.RoomChangeSkipped,
		len(summary.UnmatchedSubjectNames), len(summary.UnmatchedRoomNames),
	)
}
