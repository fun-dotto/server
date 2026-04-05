package handler

import (
	"context"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type subjectService interface {
	List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
	GetSyllabus(ctx context.Context, subjectID string) (domain.Syllabus, error)
}

type facultyService interface {
	List(ctx context.Context, ids []string) ([]domain.Faculty, error)
	GetByID(ctx context.Context, id string) (domain.Faculty, error)
	Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type roomService interface {
	List(ctx context.Context, filter domain.RoomListFilter) ([]domain.Room, error)
	GetByID(ctx context.Context, id string) (domain.Room, error)
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	Update(ctx context.Context, room domain.Room) (domain.Room, error)
	Delete(ctx context.Context, id string) error
}

type timetableItemService interface {
	List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error)
	Create(ctx context.Context, item domain.TimetableItem) (domain.TimetableItem, error)
	Delete(ctx context.Context, id string) error
}

type courseRegistrationService interface {
	List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error)
	Create(ctx context.Context, cr domain.CourseRegistration) (domain.CourseRegistration, error)
	Delete(ctx context.Context, id string) error
}

type personalCalendarItemService interface {
	List(ctx context.Context, userID string, dates []time.Time) ([]domain.PersonalCalendarItem, error)
}

type cancelledClassService interface {
	List(ctx context.Context, filter domain.CancelledClassListFilter) ([]domain.CancelledClass, error)
	GetByID(ctx context.Context, id string) (domain.CancelledClass, error)
	Create(ctx context.Context, cc domain.CancelledClass) (domain.CancelledClass, error)
	Delete(ctx context.Context, id string) error
}

type makeupClassService interface {
	List(ctx context.Context, filter domain.MakeupClassListFilter) ([]domain.MakeupClass, error)
	GetByID(ctx context.Context, id string) (domain.MakeupClass, error)
	Create(ctx context.Context, mc domain.MakeupClass) (domain.MakeupClass, error)
	Delete(ctx context.Context, id string) error
}

type roomChangeService interface {
	List(ctx context.Context, filter domain.RoomChangeListFilter) ([]domain.RoomChange, error)
	GetByID(ctx context.Context, id string) (domain.RoomChange, error)
	Create(ctx context.Context, rc domain.RoomChange) (domain.RoomChange, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	subjectSvc               subjectService
	facultySvc               facultyService
	roomSvc                  roomService
	timetableItemSvc         timetableItemService
	courseRegistrationSvc    courseRegistrationService
	personalCalendarItemSvc personalCalendarItemService
	cancelledClassSvc       cancelledClassService
	makeupClassSvc          makeupClassService
	roomChangeSvc           roomChangeService
}

func NewHandler(
	subjectSvc subjectService,
	facultySvc facultyService,
	roomSvc roomService,
	timetableItemSvc timetableItemService,
	courseRegistrationSvc courseRegistrationService,
	personalCalendarItemSvc personalCalendarItemService,
	cancelledClassSvc cancelledClassService,
	makeupClassSvc makeupClassService,
	roomChangeSvc roomChangeService,
) *Handler {
	return &Handler{
		subjectSvc:               subjectSvc,
		facultySvc:               facultySvc,
		roomSvc:                  roomSvc,
		timetableItemSvc:         timetableItemSvc,
		courseRegistrationSvc:    courseRegistrationSvc,
		personalCalendarItemSvc: personalCalendarItemSvc,
		cancelledClassSvc:       cancelledClassSvc,
		makeupClassSvc:          makeupClassSvc,
		roomChangeSvc:           roomChangeSvc,
	}
}
