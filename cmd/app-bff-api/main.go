package main

import (
	"context"
	"log"

	firebaseAdmin "firebase.google.com/go/v4"
	openapispec "github.com/fun-dotto/server/api/openapi/app"
	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/handler"
	"github.com/fun-dotto/server/internal/modules/app/middleware"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/fun-dotto/server/internal/shared/apiclient"
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
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

	spec, err := openapi3.NewLoader().LoadFromData(openapispec.Spec)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}
	spec.Servers = nil

	// 認証は AppCheckMiddleware / AuthMiddleware が担うため、バリデータでは検証しない。
	requestValidator := ginmiddleware.OapiRequestValidatorWithOptions(spec, &ginmiddleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	})

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

	funchRepository := repository.NewFunchRepository(clients.Funch)
	funchService := service.NewFunchService(funchRepository)

	h := handler.NewHandler(
		handler.WithAnnouncementService(announcementService),
		handler.WithAcademicService(academicService),
		handler.WithUserService(userService),
		handler.WithFunchService(funchService),
	)

	strictHandler := api.NewStrictHandler(h, nil)
	api.RegisterHandlersWithOptions(router, strictHandler, api.GinServerOptions{
		Middlewares: []api.MiddlewareFunc{
			api.MiddlewareFunc(middleware.AppCheckMiddleware(appCheckClient)),
			api.MiddlewareFunc(middleware.AuthMiddleware(authClient)),
			api.MiddlewareFunc(requestValidator),
		},
	})

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
