package main

import (
	"context"
	"log"
	"os"

	"github.com/fun-dotto/server/internal/modules/bus/repository"
	"github.com/fun-dotto/server/internal/modules/bus/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	url := os.Getenv("GTFS_SCHEDULE_URL")
	if url == "" {
		log.Fatal("GTFS_SCHEDULE_URL is not set")
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

	repo := repository.NewScheduleRepository(conn)
	svc := service.NewScheduleService(repo, nil)

	if err := svc.Import(context.Background(), url); err != nil {
		log.Fatalf("Failed to import gtfs: %v", err)
	}
	log.Println("gtfs import completed")
}
