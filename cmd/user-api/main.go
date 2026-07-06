package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/user/handler"
	"github.com/fun-dotto/server/internal/modules/user/middleware"
	"github.com/fun-dotto/server/internal/modules/user/openapispec"
	"github.com/fun-dotto/server/internal/modules/user/repository"
	"github.com/fun-dotto/server/internal/modules/user/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimw "github.com/oapi-codegen/gin-middleware"
)

const handlerTimeout = 15 * time.Second

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

	// マイグレーションは Atlas 専用 Cloud Run Job (cmd/migrate) で適用するため
	// API プロセス起動時の AutoMigrate 呼び出しは廃止する。

	spec, err := openapi3.NewLoader().LoadFromData(openapispec.Spec)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	spec.Servers = nil

	router := gin.Default()

	router.Use(middleware.Timeout(handlerTimeout))
	router.Use(oapimw.OapiRequestValidator(spec))

	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}
	messagingClient, err := firebaseApp.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Messaging client: %v", err)
	}

	userRepo := repository.NewUserRepository(conn)
	fcmTokenRepo := repository.NewFCMTokenRepository(conn)
	notificationRepo := repository.NewNotificationRepository(conn)
	userService := service.NewUserService(userRepo)
	fcmTokenService := service.NewFCMTokenService(fcmTokenRepo)
	notificationService := service.NewNotificationService(notificationRepo, fcmTokenRepo, messagingClient)
	h := handler.NewHandler(userService, fcmTokenService, notificationService)
	strictHandler := api.NewStrictHandler(h, []api.StrictMiddlewareFunc{
		middleware.DeadlineErrorMapper(),
	})
	api.RegisterHandlers(router, strictHandler)

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
