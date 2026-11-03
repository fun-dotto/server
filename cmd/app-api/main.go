package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	firebaseAdmin "firebase.google.com/go/v4"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/adminapi/middleware"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/modules/appapi/handler"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
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

	// Cloudflare Access を経由したリクエストは App Check の検証を省略する。
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

	calendar, err := academiccalendar.Load()
	if err != nil {
		log.Fatalf("Failed to load academic calendar: %v", err)
	}

	h := handler.New(conn, calendar)

	appCheck := auth.AppCheckTokenVerifierFunc(func(token string) error {
		_, err := appCheckClient.VerifyToken(token)
		return err
	})
	opts := connect.WithInterceptors(
		timeoutInterceptor(handlerTimeout),
		errorRecordInterceptor(),
		auth.NewInterceptor(appCheck, cfAccess, auth.FirebaseIDTokenVerifier{Client: authClient}),
	)

	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewAnnouncementServiceHandler(h, opts))
	mux.Handle(appv1connect.NewCancelledClassServiceHandler(h, opts))
	mux.Handle(appv1connect.NewCourseRegistrationServiceHandler(h, opts))
	mux.Handle(appv1connect.NewMakeupClassServiceHandler(h, opts))
	mux.Handle(appv1connect.NewPersonalCalendarItemServiceHandler(h, opts))
	mux.Handle(appv1connect.NewReservationServiceHandler(h, opts))
	mux.Handle(appv1connect.NewRoomChangeServiceHandler(h, opts))
	mux.Handle(appv1connect.NewSubjectServiceHandler(h, opts))
	mux.Handle(appv1connect.NewTimetableItemServiceHandler(h, opts))
	mux.Handle(appv1connect.NewUserServiceHandler(h, opts))
	mux.Handle(appv1connect.NewFcmTokenServiceHandler(h, opts))
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
