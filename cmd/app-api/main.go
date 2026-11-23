package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	firebaseAdmin "firebase.google.com/go/v4"
	"github.com/joho/godotenv"

	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/modules/appapi/handler"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
	"github.com/fun-dotto/server/internal/shared/cfaccess"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/interceptor"
	"github.com/fun-dotto/server/internal/shared/logging"
	"github.com/fun-dotto/server/internal/shared/microcms"
	"github.com/fun-dotto/server/internal/shared/server"
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
	// Access をかけていない環境（prd）では検証しない。
	var cfAccess auth.CFAccessVerifier
	if os.Getenv("CF_ACCESS_ENABLED") == "true" {
		v, err := cfaccess.NewVerifier(os.Getenv("CF_ACCESS_ISSUER"), os.Getenv("CF_ACCESS_AUD"))
		if err != nil {
			log.Fatalf("Failed to initialize Cloudflare Access verifier: %v", err)
		}
		cfAccess = v
	}

	// 許可されたサービス（Web）からの Google ID トークン付きリクエストは App Check の検証を省略する。
	service, err := auth.NewGoogleServiceVerifier(os.Getenv("ID_TOKEN_AUDIENCE"), os.Getenv("ID_TOKEN_ALLOWED_SERVICE_ACCOUNTS"))
	if err != nil {
		log.Fatalf("Failed to initialize service ID token verifier: %v", err)
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

	// microCMS が未設定でも App API 全体は起動させ、ListArticles だけ UNAVAILABLE にする。
	var articles handler.ArticleLister
	var articleBaseURL string
	if c, err := microcms.NewClientFromEnv(); err != nil {
		log.Printf("Warning: articles are unavailable: %v", err)
	} else if u, err := microcms.ArticleBaseURLFromEnv(); err != nil {
		log.Printf("Warning: articles are unavailable: %v", err)
	} else {
		articles, articleBaseURL = c, u
	}

	h := handler.New(conn, calendar, articles, articleBaseURL)

	appCheck := auth.AppCheckTokenVerifierFunc(func(token string) error {
		_, err := appCheckClient.VerifyToken(token)
		return err
	})
	opts := connect.WithInterceptors(
		interceptor.Timeout(handlerTimeout),
		interceptor.RecordError(),
		auth.NewInterceptor(appCheck, cfAccess, service, auth.FirebaseIDTokenVerifier{Client: authClient}),
	)

	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewArticleServiceHandler(h, opts))
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
	// バスの運行情報・学食のメニューは未実装（UNIMPLEMENTED を返す）。
	mux.Handle(appv1connect.NewBusTripServiceHandler(appv1connect.UnimplementedBusTripServiceHandler{}, opts))
	mux.Handle(appv1connect.NewBusTimetableStopServiceHandler(appv1connect.UnimplementedBusTimetableStopServiceHandler{}, opts))
	mux.Handle(appv1connect.NewMenuItemServiceHandler(appv1connect.UnimplementedMenuItemServiceHandler{}, opts))

	if err := server.Run(logging.Middleware(mux), server.Addr()); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
