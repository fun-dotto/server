// Package server provides the common HTTP server bootstrap (timeouts,
// graceful shutdown) shared by every API cmd (academic-api, user-api,
// announcement-api, ...).
package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 8 * time.Second

	defaultPort = "8080"
)

// Addr returns the address the HTTP server should listen on. 環境変数 PORT が
// 設定されていればそれを使う (Cloud Run が本番で注入し、ローカルでは portless が
// 4000-4999 の空きポートを注入する)。未設定なら :8080 にフォールバックする。
func Addr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}
	return ":" + port
}

// Run starts an HTTP server on addr with the project's standard timeouts,
// and blocks until it exits, either because ListenAndServe failed or
// because SIGINT/SIGTERM triggered a graceful shutdown.
func Run(handler http.Handler, addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
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
			return err
		}
		return nil
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
		return err
	}
	return nil
}
