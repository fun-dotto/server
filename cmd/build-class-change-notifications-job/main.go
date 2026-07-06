package main

import (
	"context"
	"log"

	academicrepository "github.com/fun-dotto/server/internal/modules/academic/repository"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	userrepository "github.com/fun-dotto/server/internal/modules/user/repository"
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

	cancelledRepo := academicrepository.NewCancelledClassRepository(conn)
	makeupRepo := academicrepository.NewMakeupClassRepository(conn)
	roomChangeRepo := academicrepository.NewRoomChangeRepository(conn)
	courseRegRepo := academicrepository.NewCourseRegistrationRepository(conn)
	notificationRepo := userrepository.NewNotificationRepository(conn)

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
