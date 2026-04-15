package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fun-dotto/academic-api/assets"
	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/event"
	"github.com/fun-dotto/academic-api/internal/handler"
	"github.com/fun-dotto/academic-api/internal/middleware"
	"github.com/fun-dotto/academic-api/internal/repository"
	"github.com/fun-dotto/academic-api/internal/service"
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

	// Repositories
	subjectRepo := repository.NewSubjectRepository(db)
	syllabusRepo := repository.NewSyllabusRepository(db)
	facultyRepo := repository.NewFacultyRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	timetableItemRepo := repository.NewTimetableItemRepository(db)
	courseRegistrationRepo := repository.NewCourseRegistrationRepository(db)
	cancelledClassRepo := repository.NewCancelledClassRepository(db)
	makeupClassRepo := repository.NewMakeupClassRepository(db)
	roomChangeRepo := repository.NewRoomChangeRepository(db)
	// Events
	substituteDayMap, err := event.LoadSubstituteDayMap(assets.EventsJSON)
	if err != nil {
		log.Fatalf("Failed to load substitute day map: %v", err)
	}

	// Services
	subjectSvc := service.NewSubjectService(subjectRepo, syllabusRepo)
	facultySvc := service.NewFacultyService(facultyRepo)
	roomSvc := service.NewRoomService(roomRepo)
	timetableItemSvc := service.NewTimetableItemService(timetableItemRepo)
	courseRegistrationSvc := service.NewCourseRegistrationService(courseRegistrationRepo)
	personalCalendarItemSvc := service.NewPersonalCalendarItemService(courseRegistrationRepo, timetableItemRepo, cancelledClassRepo, makeupClassRepo, roomChangeRepo, substituteDayMap)
	cancelledClassSvc := service.NewCancelledClassService(cancelledClassRepo)
	makeupClassSvc := service.NewMakeupClassService(makeupClassRepo)
	roomChangeSvc := service.NewRoomChangeService(roomChangeRepo)

	// Handler + Router
	h := handler.NewHandler(subjectSvc, facultySvc, roomSvc, timetableItemSvc, courseRegistrationSvc, personalCalendarItemSvc, cancelledClassSvc, makeupClassSvc, roomChangeSvc)
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
