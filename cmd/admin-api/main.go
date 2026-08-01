package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/handler"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
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
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to get Firebase Auth client: %v", err)
	}

	// 仕様の原本は api/openapi/admin/openapi.yaml だけに置く。runtime image にも
	// 同梱しているため、WORKDIR "/" 基準でこの相対パスに解決される。
	spec, err := openapi3.NewLoader().LoadFromFile("api/openapi/admin/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	spec.Servers = nil

	router := gin.Default()

	router.Use(ginmiddleware.OapiRequestValidatorWithOptions(spec, &ginmiddleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			if authStatusCode, authMessage, ok := middleware.GetAuthenticationError(c); ok {
				c.AbortWithStatusJSON(authStatusCode, gin.H{"error": authMessage})
				return
			}
			c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: middleware.FirebaseAuthenticationFunc(authClient),
		},
	}))

	clients, err := apiclient.NewExternalClients(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize external clients: %v", err)
	}

	h := handler.NewHandler(clients.Academic, clients.Announcement, clients.Funch, clients.User)
	api.RegisterHandlers(router, h)

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
