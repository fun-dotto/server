package main

import (
	"log"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/database"
	"github.com/fun-dotto/announcement-api/internal/handler"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	db, err := database.Connect()
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

	router := gin.Default()

	announcementRepository := repository.NewAnnouncementRepository(db)

	announcementService := service.NewAnnouncementService(announcementRepository)

	h := handler.NewHandler(announcementService)

	api.RegisterHandlers(router, h)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
