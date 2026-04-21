package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	api "github.com/fun-dotto/user-api/generated"
	"github.com/fun-dotto/user-api/internal/database"
	"github.com/fun-dotto/user-api/internal/handler"
	"github.com/fun-dotto/user-api/internal/middleware"
	"github.com/fun-dotto/user-api/internal/repository"
	"github.com/fun-dotto/user-api/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimw "github.com/oapi-codegen/gin-middleware"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
	handlerTimeout    = 15 * time.Second
	shutdownTimeout   = 8 * time.Second
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	db, err := database.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	spec, err := openapi3.NewLoader().LoadFromFile("openapi/openapi.yaml")
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

	userRepo := repository.NewUserRepository(db)
	fcmTokenRepo := repository.NewFCMTokenRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	userService := service.NewUserService(userRepo)
	fcmTokenService := service.NewFCMTokenService(fcmTokenRepo)
	notificationService := service.NewNotificationService(notificationRepo, fcmTokenRepo, messagingClient)
	h := handler.NewHandler(userService, fcmTokenService, notificationService)
	strictHandler := api.NewStrictHandler(h, []api.StrictMiddlewareFunc{
		middleware.DeadlineErrorMapper(),
	})
	api.RegisterHandlers(router, strictHandler)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to start server: %v", err)
		}
		return
	case <-ctx.Done():
		log.Println("Shutdown signal received, draining in-flight requests...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		if closeErr := srv.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
			log.Printf("Server force close error: %v", closeErr)
		}
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server exited with error: %v", err)
	}
}
