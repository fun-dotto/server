package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/gen/admin/v1/adminv1connect"
	"github.com/fun-dotto/server/internal/modules/adminapi/handler"
	"github.com/fun-dotto/server/internal/modules/adminapi/repository"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/joho/godotenv"
)

const handlerTimeout = 15 * time.Second

func main() {
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

	opts := connect.WithInterceptors(timeoutInterceptor(handlerTimeout), errorRecordInterceptor())

	mux := http.NewServeMux()
	mux.Handle(adminv1connect.NewAnnouncementServiceHandler(handler.NewAnnouncementHandler(repository.NewAnnouncementRepository(conn)), opts))
	mux.Handle(adminv1connect.NewCalendarServiceHandler(handler.NewCalendarHandler(repository.NewCalendarRepository(conn)), opts))
	mux.Handle(adminv1connect.NewCalendarDateServiceHandler(handler.NewCalendarDateHandler(repository.NewCalendarDateRepository(conn)), opts))
	mux.Handle(adminv1connect.NewCancelledClassServiceHandler(handler.NewCancelledClassHandler(repository.NewCancelledClassRepository(conn)), opts))
	mux.Handle(adminv1connect.NewCourseRegistrationServiceHandler(handler.NewCourseRegistrationHandler(repository.NewCourseRegistrationRepository(conn)), opts))
	mux.Handle(adminv1connect.NewFacultyServiceHandler(handler.NewFacultyHandler(repository.NewFacultyRepository(conn)), opts))
	mux.Handle(adminv1connect.NewFacultyRoomServiceHandler(handler.NewFacultyRoomHandler(repository.NewFacultyRoomRepository(conn)), opts))
	mux.Handle(adminv1connect.NewFareAttributeServiceHandler(handler.NewFareAttributeHandler(repository.NewFareAttributeRepository(conn)), opts))
	mux.Handle(adminv1connect.NewFareRuleServiceHandler(handler.NewFareRuleHandler(repository.NewFareRuleRepository(conn)), opts))
	mux.Handle(adminv1connect.NewFcmTokenServiceHandler(handler.NewFcmTokenHandler(repository.NewFCMTokenRepository(conn)), opts))
	mux.Handle(adminv1connect.NewMakeupClassServiceHandler(handler.NewMakeupClassHandler(repository.NewMakeupClassRepository(conn)), opts))
	mux.Handle(adminv1connect.NewNotificationServiceHandler(handler.NewNotificationHandler(repository.NewNotificationRepository(conn)), opts))
	mux.Handle(adminv1connect.NewNotificationTargetUserServiceHandler(handler.NewNotificationTargetUserHandler(repository.NewNotificationTargetUserRepository(conn)), opts))
	mux.Handle(adminv1connect.NewRoomServiceHandler(handler.NewRoomHandler(repository.NewRoomRepository(conn)), opts))
	mux.Handle(adminv1connect.NewRoomChangeServiceHandler(handler.NewRoomChangeHandler(repository.NewRoomChangeRepository(conn)), opts))
	mux.Handle(adminv1connect.NewRoomReservationServiceHandler(handler.NewRoomReservationHandler(repository.NewRoomReservationRepository(conn)), opts))
	mux.Handle(adminv1connect.NewRouteServiceHandler(handler.NewRouteHandler(repository.NewRouteRepository(conn)), opts))
	mux.Handle(adminv1connect.NewStopServiceHandler(handler.NewStopHandler(repository.NewStopRepository(conn)), opts))
	mux.Handle(adminv1connect.NewStopTimeServiceHandler(handler.NewStopTimeHandler(repository.NewStopTimeRepository(conn)), opts))
	mux.Handle(adminv1connect.NewSubjectServiceHandler(handler.NewSubjectHandler(repository.NewSubjectRepository(conn)), opts))
	mux.Handle(adminv1connect.NewSyllabusServiceHandler(handler.NewSyllabusHandler(repository.NewSyllabusRepository(conn)), opts))
	mux.Handle(adminv1connect.NewTimetableItemServiceHandler(handler.NewTimetableItemHandler(repository.NewTimetableItemRepository(conn)), opts))
	mux.Handle(adminv1connect.NewTripServiceHandler(handler.NewTripHandler(repository.NewTripRepository(conn)), opts))
	mux.Handle(adminv1connect.NewUserServiceHandler(handler.NewUserHandler(repository.NewUserRepository(conn)), opts))
	mux.Handle(adminv1connect.NewZoneServiceHandler(handler.NewZoneHandler(repository.NewZoneRepository(conn)), opts))

	if err := server.Run(logging.Middleware(mux), server.Addr()); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}

// timeoutInterceptor は 1 リクエストあたりの処理時間の上限を設ける。
func timeoutInterceptor(d time.Duration) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx, req)
		}
	}
}

// errorRecordInterceptor は RPC のエラーをアクセスログに記録する。
func errorRecordInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				logging.RecordError(ctx, err)
			}
			return res, err
		}
	}
}
