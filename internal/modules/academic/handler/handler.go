package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
	"github.com/fun-dotto/subject-api/internal/domain"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type courseService interface {
	List(ctx context.Context) ([]domain.Course, error)
	GetByID(ctx context.Context, id string) (domain.Course, error)
	Create(ctx context.Context, course domain.Course) (domain.Course, error)
	Update(ctx context.Context, id string, course domain.Course) (domain.Course, error)
	Delete(ctx context.Context, id string) error
}

type facultyService interface {
	List(ctx context.Context) ([]domain.Faculty, error)
	GetByID(ctx context.Context, id string) (domain.Faculty, error)
	Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error)
	Update(ctx context.Context, id string, faculty domain.Faculty) (domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type dayOfWeekTimetableSlotService interface {
	List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error)
	GetByID(ctx context.Context, id string) (domain.DayOfWeekTimetableSlot, error)
	Create(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error)
	Update(ctx context.Context, id string, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error)
	Delete(ctx context.Context, id string) error
}

type subjectCategoryService interface {
	List(ctx context.Context) ([]domain.SubjectCategory, error)
	GetByID(ctx context.Context, id string) (domain.SubjectCategory, error)
	Create(ctx context.Context, category domain.SubjectCategory) (domain.SubjectCategory, error)
	Update(ctx context.Context, id string, category domain.SubjectCategory) (domain.SubjectCategory, error)
	Delete(ctx context.Context, id string) error
}

type subjectService interface {
	List(ctx context.Context) ([]domain.Subject, error)
	GetByID(ctx context.Context, id string) (domain.Subject, error)
	Create(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	Update(ctx context.Context, id string, subject domain.Subject) (domain.Subject, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	courseSvc   courseService
	facultySvc  facultyService
	slotSvc     dayOfWeekTimetableSlotService
	categorySvc subjectCategoryService
	subjectSvc  subjectService
}

func NewHandler(
	courseSvc courseService,
	facultySvc facultyService,
	slotSvc dayOfWeekTimetableSlotService,
	categorySvc subjectCategoryService,
	subjectSvc subjectService,
) *Handler {
	return &Handler{
		courseSvc:   courseSvc,
		facultySvc:  facultySvc,
		slotSvc:     slotSvc,
		categorySvc: categorySvc,
		subjectSvc:  subjectSvc,
	}
}
