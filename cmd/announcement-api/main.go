package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	api "github.com/fun-dotto/server/gen/announcement"
	"github.com/fun-dotto/server/internal/modules/announcement/handler"
	"github.com/fun-dotto/server/internal/modules/announcement/openapispec"
	"github.com/fun-dotto/server/internal/modules/announcement/repository"
	"github.com/fun-dotto/server/internal/modules/announcement/service"
	"github.com/fun-dotto/server/internal/shared/db"
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
	shutdownTimeout   = 8 * time.Second
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

	// マイグレーションは Atlas 専用 Cloud Run Job (cmd/migrate) で適用するため
	// API プロセス起動時の AutoMigrate 呼び出しは行わない。

	spec, err := openapi3.NewLoader().LoadFromData(openapispec.Spec)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	spec.Servers = nil

	router := gin.Default()

	router.Use(oapimw.OapiRequestValidator(spec))

	// Repositories
	announcementRepo := repository.NewAnnouncementRepository(conn)

	// Services
	announcementSvc := service.NewAnnouncementService(announcementRepo)

	// Handler + Router
	h := handler.NewHandler(announcementSvc)
	strictHandler := api.NewStrictHandler(h, nil)
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
