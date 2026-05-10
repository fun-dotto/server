package main

import (
	"context"
	"log"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/repository"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
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

	cancelledRepo := repository.NewCancelledClassRepository(conn)
	makeupRepo := repository.NewMakeupClassRepository(conn)
	roomChangeRepo := repository.NewRoomChangeRepository(conn)
	courseRegRepo := repository.NewCourseRegistrationRepository(conn)
	notificationRepo := repository.NewNotificationRepository(conn)

	svc := service.NewClassChangeNotificationService(
		cancelledRepo,
		makeupRepo,
		roomChangeRepo,
		courseRegRepo,
		notificationRepo,
	)

	summary, err := svc.EnqueueNotifications(context.Background())
	if err != nil {
		log.Fatalf("Failed to enqueue notifications: %v", err)
	}
	log.Printf(
		"enqueue summary: cancelled=%d makeup=%d room_change=%d skipped=%d",
		summary.CancelledEnqueued,
		summary.MakeupEnqueued,
		summary.RoomChangeEnqueued,
		summary.Skipped,
	)
}
