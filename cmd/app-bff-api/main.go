package main

import (
	"context"
	"log"

	firebaseAdmin "firebase.google.com/go/v4"
	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/handler"
	"github.com/fun-dotto/server/internal/modules/app/middleware"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/fun-dotto/server/internal/shared/apiclient"
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	// Firebase App Check の初期化
	app, err := firebaseAdmin.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	appCheckClient, err := app.AppCheck(ctx)
	if err != nil {
		log.Fatalf("error initializing App Check client: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing Auth client: %v\n", err)
	}

	router := gin.Default()

	// 外部APIクライアントを初期化
	clients, err := apiclient.NewExternalClients(ctx)
	if err != nil {
		log.Fatalf("Failed to create external clients: %v", err)
	}

	announcementRepository := repository.NewAnnouncementRepository(clients.Announcement)
	announcementService := service.NewAnnouncementService(announcementRepository)

	academicRepository := repository.NewAcademicRepository(clients.Academic)
	academicService := service.NewAcademicService(academicRepository)

	userRepository := repository.NewUserRepository(clients.User)
	userService := service.NewUserService(userRepository)

	opts := []handler.Option{
		handler.WithAnnouncementService(announcementService),
		handler.WithAcademicService(academicService),
		handler.WithUserService(userService),
	}
	if clients.Funch != nil {
		funchRepository := repository.NewFunchRepository(clients.Funch)
		opts = append(opts, handler.WithFunchService(service.NewFunchService(funchRepository)))
	} else {
		log.Printf("Warning: FUNCH_API_URL is not set; funch endpoints are disabled")
	}

	h := handler.NewHandler(opts...)

	strictHandler := api.NewStrictHandler(h, nil)
	api.RegisterHandlersWithOptions(router, strictHandler, api.GinServerOptions{
		Middlewares: []api.MiddlewareFunc{
			api.MiddlewareFunc(middleware.AppCheckMiddleware(appCheckClient)),
			api.MiddlewareFunc(middleware.AuthMiddleware(authClient)),
		},
	})

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
