package main

import (
	"errors"
	"log"
	"net/http"
	"time"

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

	userRepo := repository.NewUserRepository(db)
	fcmTokenRepo := repository.NewFCMTokenRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	userService := service.NewUserService(userRepo)
	fcmTokenService := service.NewFCMTokenService(fcmTokenRepo)
	notificationService := service.NewNotificationService(notificationRepo)
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
	log.Printf("Server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Failed to start server:", err)
	}
}
