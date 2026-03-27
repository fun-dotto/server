package main

import (
	"log"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/handler"
	"github.com/fun-dotto/academic-api/internal/repository"
	"github.com/fun-dotto/academic-api/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	middleware "github.com/oapi-codegen/gin-middleware"
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

	router.Use(middleware.OapiRequestValidator(spec))

	// Repositories
	subjectRepo := repository.NewSubjectRepository(db)
	syllabusRepo := repository.NewSyllabusRepository(db)
	facultyRepo := repository.NewFacultyRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	timetableItemRepo := repository.NewTimetableItemRepository(db)
	// Services
	subjectSvc := service.NewSubjectService(subjectRepo, syllabusRepo)
	facultySvc := service.NewFacultyService(facultyRepo)
	roomSvc := service.NewRoomService(roomRepo)
	timetableItemSvc := service.NewTimetableItemService(timetableItemRepo)

	// Handler + Router
	h := handler.NewHandler(subjectSvc, facultySvc, roomSvc, timetableItemSvc)
	strictHandler := api.NewStrictHandler(h, nil)
	api.RegisterHandlers(router, strictHandler)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
