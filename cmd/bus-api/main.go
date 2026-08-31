package main

import (
	"log"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/handler"
	"github.com/fun-dotto/server/internal/modules/bus/openapispec"
	"github.com/fun-dotto/server/internal/modules/bus/repository"
	"github.com/fun-dotto/server/internal/modules/bus/service"
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

	spec, err := openapi3.NewLoader().LoadFromData(openapispec.Spec)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}
	spec.Servers = nil

	router := gin.Default()
	router.Use(oapimw.OapiRequestValidator(spec))

	busRepo := repository.NewBusRepository(conn)
	busSvc := service.NewBusService(busRepo)
	h := handler.NewHandler(busSvc)
	strictHandler := api.NewStrictHandler(h, nil)
	api.RegisterHandlers(router, strictHandler)

	if err := server.Run(router, server.Addr()); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
