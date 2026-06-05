package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/assets"
	"github.com/fun-dotto/server/internal/modules/academic/event"
	"github.com/fun-dotto/server/internal/modules/academic/handler"
	"github.com/fun-dotto/server/internal/modules/academic/middleware"
	"github.com/fun-dotto/server/internal/modules/academic/openapispec"
	"github.com/fun-dotto/server/internal/modules/academic/repository"
	"github.com/fun-dotto/server/internal/modules/academic/service"
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
	handlerTimeout    = 15 * time.Second
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

	// マイグレーションは Atlas 専用 Cloud Run Job (cmd/migrate-job) で適用するため
	// API プロセス起動時の AutoMigrate 呼び出しは廃止する (Notion 計画 §3 / §7-C)。

	spec, err := openapi3.NewLoader().LoadFromData(openapispec.Spec)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	spec.Servers = nil

	router := gin.Default()

	router.Use(middleware.Timeout(handlerTimeout))
	router.Use(oapimw.OapiRequestValidator(spec))

	// Repositories
	subjectRepo := repository.NewSubjectRepository(conn)
	syllabusRepo := repository.NewSyllabusRepository(conn)
	facultyRepo := repository.NewFacultyRepository(conn)
	roomRepo := repository.NewRoomRepository(conn)
	timetableItemRepo := repository.NewTimetableItemRepository(conn)
	courseRegistrationRepo := repository.NewCourseRegistrationRepository(conn)
	cancelledClassRepo := repository.NewCancelledClassRepository(conn)
	makeupClassRepo := repository.NewMakeupClassRepository(conn)
	roomChangeRepo := repository.NewRoomChangeRepository(conn)
	facultyRoomRepo := repository.NewFacultyRoomRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	// Events
	substituteDayMap, err := event.LoadSubstituteDayMap(assets.EventsJSON)
	if err != nil {
		log.Fatalf("Failed to load substitute day map: %v", err)
	}

	holidaySet, err := event.LoadHolidaySet(assets.HolidaysJSON)
	if err != nil {
		log.Fatalf("Failed to load holiday set: %v", err)
	}

	// Services
	subjectSvc := service.NewSubjectService(subjectRepo, syllabusRepo)
	facultySvc := service.NewFacultyService(facultyRepo)
	roomSvc := service.NewRoomService(roomRepo)
	timetableItemSvc := service.NewTimetableItemService(timetableItemRepo)
	courseRegistrationSvc := service.NewCourseRegistrationService(courseRegistrationRepo)
	personalCalendarItemSvc := service.NewPersonalCalendarItemService(courseRegistrationRepo, timetableItemRepo, cancelledClassRepo, makeupClassRepo, roomChangeRepo, substituteDayMap, holidaySet)
	cancelledClassSvc := service.NewCancelledClassService(cancelledClassRepo)
	makeupClassSvc := service.NewMakeupClassService(makeupClassRepo)
	roomChangeSvc := service.NewRoomChangeService(roomChangeRepo)
	facultyRoomSvc := service.NewFacultyRoomService(facultyRoomRepo)

	// Handler + Router
	h := handler.NewHandler(subjectSvc, facultySvc, roomSvc, timetableItemSvc, courseRegistrationSvc, personalCalendarItemSvc, cancelledClassSvc, makeupClassSvc, roomChangeSvc, facultyRoomSvc, userRepo)
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
