package main

import (
	"log"

	api "github.com/fun-dotto/server/gen/announcement"
	"github.com/fun-dotto/server/internal/modules/announcement/handler"
	"github.com/fun-dotto/server/internal/modules/announcement/openapispec"
	"github.com/fun-dotto/server/internal/modules/announcement/repository"
	"github.com/fun-dotto/server/internal/modules/announcement/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimw "github.com/oapi-codegen/gin-middleware"
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

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
