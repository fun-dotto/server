package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/gtfs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	scheduleURL := os.Getenv("GTFS_SCHEDULE_URL")
	if scheduleURL == "" {
		return fmt.Errorf("GTFS_SCHEDULE_URL is not set")
	}

	conn, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	return service.New(conn).Import(context.Background(), scheduleURL)
}
