package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/gen/admin/v1/adminv1connect"
	"github.com/fun-dotto/server/internal/modules/adminapi/handler"
	"github.com/fun-dotto/server/internal/modules/adminapi/middleware"
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

	// Cloudflare Access を経由していないリクエストは拒否する。
	cfAccess, err := middleware.NewCFAccessVerifier(os.Getenv("CF_ACCESS_ISSUER"), os.Getenv("CF_ACCESS_AUD"))
	if err != nil {
		log.Fatalf("Failed to initialize Cloudflare Access verifier: %v", err)
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
	mux.Handle(adminv1connect.NewAnnouncementServiceHandler(handler.NewAnnouncementHandler(conn), opts))
	mux.Handle(adminv1connect.NewCalendarServiceHandler(handler.NewCalendarHandler(conn), opts))
	mux.Handle(adminv1connect.NewCalendarDateServiceHandler(handler.NewCalendarDateHandler(conn), opts))
	mux.Handle(adminv1connect.NewCancelledClassServiceHandler(handler.NewCancelledClassHandler(conn), opts))
	mux.Handle(adminv1connect.NewCourseRegistrationServiceHandler(handler.NewCourseRegistrationHandler(conn), opts))
	mux.Handle(adminv1connect.NewFacultyServiceHandler(handler.NewFacultyHandler(conn), opts))
	mux.Handle(adminv1connect.NewFacultyRoomServiceHandler(handler.NewFacultyRoomHandler(conn), opts))
	mux.Handle(adminv1connect.NewFareAttributeServiceHandler(handler.NewFareAttributeHandler(conn), opts))
	mux.Handle(adminv1connect.NewFareRuleServiceHandler(handler.NewFareRuleHandler(conn), opts))
	mux.Handle(adminv1connect.NewFcmTokenServiceHandler(handler.NewFcmTokenHandler(conn), opts))
	mux.Handle(adminv1connect.NewMakeupClassServiceHandler(handler.NewMakeupClassHandler(conn), opts))
	mux.Handle(adminv1connect.NewNotificationServiceHandler(handler.NewNotificationHandler(conn), opts))
	mux.Handle(adminv1connect.NewNotificationTargetUserServiceHandler(handler.NewNotificationTargetUserHandler(conn), opts))
	mux.Handle(adminv1connect.NewRoomServiceHandler(handler.NewRoomHandler(conn), opts))
	mux.Handle(adminv1connect.NewRoomChangeServiceHandler(handler.NewRoomChangeHandler(conn), opts))
	mux.Handle(adminv1connect.NewRoomReservationServiceHandler(handler.NewRoomReservationHandler(conn), opts))
	mux.Handle(adminv1connect.NewRouteServiceHandler(handler.NewRouteHandler(conn), opts))
	mux.Handle(adminv1connect.NewStopServiceHandler(handler.NewStopHandler(conn), opts))
	mux.Handle(adminv1connect.NewStopTimeServiceHandler(handler.NewStopTimeHandler(conn), opts))
	mux.Handle(adminv1connect.NewSubjectServiceHandler(handler.NewSubjectHandler(conn), opts))
	mux.Handle(adminv1connect.NewSyllabusServiceHandler(handler.NewSyllabusHandler(conn), opts))
	mux.Handle(adminv1connect.NewTimetableItemServiceHandler(handler.NewTimetableItemHandler(conn), opts))
	mux.Handle(adminv1connect.NewTripServiceHandler(handler.NewTripHandler(conn), opts))
	mux.Handle(adminv1connect.NewUserServiceHandler(handler.NewUserHandler(conn), opts))
	mux.Handle(adminv1connect.NewZoneServiceHandler(handler.NewZoneHandler(conn), opts))

	if err := server.Run(logging.Middleware(cfAccess.Middleware(mux)), server.Addr()); err != nil {
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
