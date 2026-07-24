package main

import (
	"log"
	"net/http"
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
	"github.com/fun-dotto/server/internal/shared/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimw "github.com/oapi-codegen/gin-middleware"
)

const handlerTimeout = 15 * time.Second

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
	userSvc := service.NewUserService(userRepo)

	// Handler + Router
	h := handler.NewHandler(subjectSvc, facultySvc, roomSvc, timetableItemSvc, courseRegistrationSvc, personalCalendarItemSvc, cancelledClassSvc, makeupClassSvc, roomChangeSvc, facultyRoomSvc, userSvc)
	strictHandler := api.NewStrictHandlerWithOptions(h, []api.StrictMiddlewareFunc{
		middleware.DeadlineErrorMapper(),
	}, api.StrictGinServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		HandlerErrorFunc:         handlerErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	})
	api.RegisterHandlers(router, strictHandler)

	if err := server.Run(router, ":8080"); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}

// oapi-codegen が生成する strict handler の既定のエラーハンドラは err.Error() を
// そのままレスポンス本文に載せるため、SQL 文やドライバのエラーメッセージといった
// 内部実装の詳細がクライアントへ露出し得る。本文は固定文言に差し替え、
// 詳細はサーバー側のログにのみ出力する。

const (
	badRequestMessage    = "invalid request"
	internalErrorMessage = "internal server error"
)

// requestErrorHandler はリクエストのパース・デコードに失敗した場合に 400 を返す。
func requestErrorHandler(c *gin.Context, err error) {
	logStrictError(c, "request error", err)
	c.JSON(http.StatusBadRequest, gin.H{"msg": badRequestMessage})
}

// handlerErrorHandler はハンドラ（および strict middleware）が non-nil error を
// 返した場合に 500 を返す。
func handlerErrorHandler(c *gin.Context, err error) {
	logStrictError(c, "handler error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"msg": internalErrorMessage})
}

// responseErrorHandler はレスポンスのシリアライズに失敗した場合、あるいは想定外の
// レスポンス型が返された場合に 500 を返す。
func responseErrorHandler(c *gin.Context, err error) {
	logStrictError(c, "response error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"msg": internalErrorMessage})
}

func logStrictError(c *gin.Context, kind string, err error) {
	log.Printf("%s: %s %s: %v", kind, c.Request.Method, c.Request.URL.Path, err)
}
