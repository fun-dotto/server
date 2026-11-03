package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	firebaseAdmin "firebase.google.com/go/v4"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/modules/appapi/handler"
	"github.com/fun-dotto/server/internal/shared/apiclient"
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

	ctx := context.Background()

	firebaseApp, err := firebaseAdmin.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}
	appCheckClient, err := firebaseApp.AppCheck(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize App Check client: %v", err)
	}
	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Auth client: %v", err)
	}

	clients, err := apiclient.NewExternalClients(ctx)
	if err != nil {
		log.Fatalf("Failed to create external clients: %v", err)
	}

	announcementService := service.NewAnnouncementService(repository.NewAnnouncementRepository(clients.Announcement))
	academicService := service.NewAcademicService(repository.NewAcademicRepository(clients.Academic))
	userService := service.NewUserService(repository.NewUserRepository(clients.User))

	appCheck := auth.AppCheckTokenVerifierFunc(func(token string) error {
		_, err := appCheckClient.VerifyToken(token)
		return err
	})
	opts := connect.WithInterceptors(
		timeoutInterceptor(handlerTimeout),
		errorRecordInterceptor(),
		auth.NewInterceptor(appCheck, auth.FirebaseIDTokenVerifier{Client: authClient}),
	)

	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewAnnouncementServiceHandler(handler.NewAnnouncementHandler(announcementService), opts))
	mux.Handle(appv1connect.NewCancelledClassServiceHandler(handler.NewCancelledClassHandler(academicService), opts))
	mux.Handle(appv1connect.NewCourseRegistrationServiceHandler(handler.NewCourseRegistrationHandler(academicService), opts))
	mux.Handle(appv1connect.NewMakeupClassServiceHandler(handler.NewMakeupClassHandler(academicService), opts))
	mux.Handle(appv1connect.NewPersonalCalendarItemServiceHandler(handler.NewPersonalCalendarItemHandler(academicService), opts))
	mux.Handle(appv1connect.NewReservationServiceHandler(handler.NewReservationHandler(academicService), opts))
	mux.Handle(appv1connect.NewRoomChangeServiceHandler(handler.NewRoomChangeHandler(academicService), opts))
	mux.Handle(appv1connect.NewSubjectServiceHandler(handler.NewSubjectHandler(academicService), opts))
	mux.Handle(appv1connect.NewTimetableItemServiceHandler(handler.NewTimetableItemHandler(academicService), opts))
	mux.Handle(appv1connect.NewUserServiceHandler(handler.NewUserHandler(userService), opts))
	mux.Handle(appv1connect.NewFcmTokenServiceHandler(handler.NewFcmTokenHandler(userService), opts))
	// バスの運行情報・学食のメニューは app-bff-api と同じく未実装（UNIMPLEMENTED を返す）。
	mux.Handle(appv1connect.NewBusTripServiceHandler(appv1connect.UnimplementedBusTripServiceHandler{}, opts))
	mux.Handle(appv1connect.NewBusTimetableStopServiceHandler(appv1connect.UnimplementedBusTimetableStopServiceHandler{}, opts))
	mux.Handle(appv1connect.NewMenuItemServiceHandler(appv1connect.UnimplementedMenuItemServiceHandler{}, opts))

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
