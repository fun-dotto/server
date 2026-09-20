package service

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
)

type AcademicRepository interface {
	GetFaculties() ([]domain.Faculty, error)
	GetFaculty(id string) (*domain.Faculty, error)
	GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error)
	GetSubject(id string) (*domain.Subject, error)
	GetCourseRegistrations(userID string, semesters []domain.CourseSemester, year *int) ([]domain.CourseRegistration, error)
	CreateCourseRegistration(userID string, subjectID string) (*domain.CourseRegistration, error)
	DeleteCourseRegistration(id string) error
	GetTimetableItems(query domain.TimetableItemQuery) ([]domain.TimetableItem, error)
	GetPersonalCalendarItems(userID string, dates []time.Time) ([]domain.PersonalCalendarItem, error)
	GetCancelledClasses(query domain.CancelledClassQuery) ([]domain.CancelledClass, error)
	GetMakeupClasses(query domain.MakeupClassQuery) ([]domain.MakeupClass, error)
	GetRoomChanges(query domain.RoomChangeQuery) ([]domain.RoomChange, error)
	GetReservations(query domain.ReservationQuery) ([]domain.Reservation, error)
}

type AcademicService struct {
	repository AcademicRepository
}

func NewAcademicService(repository AcademicRepository) *AcademicService {
	return &AcademicService{repository: repository}
}

func (s *AcademicService) GetFaculties() ([]domain.Faculty, error) {
	return s.repository.GetFaculties()
}

func (s *AcademicService) GetFaculty(id string) (*domain.Faculty, error) {
	return s.repository.GetFaculty(id)
}

func (s *AcademicService) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	return s.repository.GetSubjects(query)
}

func (s *AcademicService) GetSubject(id string) (*domain.Subject, error) {
	return s.repository.GetSubject(id)
}

func (s *AcademicService) GetCourseRegistrations(userID string, semesters []domain.CourseSemester, year *int) ([]domain.CourseRegistration, error) {
	return s.repository.GetCourseRegistrations(userID, semesters, year)
}

func (s *AcademicService) CreateCourseRegistration(userID string, subjectID string) (*domain.CourseRegistration, error) {
	return s.repository.CreateCourseRegistration(userID, subjectID)
}

func (s *AcademicService) DeleteCourseRegistration(id string) error {
	return s.repository.DeleteCourseRegistration(id)
}

func (s *AcademicService) GetTimetableItems(query domain.TimetableItemQuery) ([]domain.TimetableItem, error) {
	return s.repository.GetTimetableItems(query)
}

func (s *AcademicService) GetPersonalCalendarItems(userID string, dates []time.Time) ([]domain.PersonalCalendarItem, error) {
	return s.repository.GetPersonalCalendarItems(userID, dates)
}

func (s *AcademicService) GetCancelledClasses(query domain.CancelledClassQuery) ([]domain.CancelledClass, error) {
	return s.repository.GetCancelledClasses(query)
}

func (s *AcademicService) GetMakeupClasses(query domain.MakeupClassQuery) ([]domain.MakeupClass, error) {
	return s.repository.GetMakeupClasses(query)
}

func (s *AcademicService) GetRoomChanges(query domain.RoomChangeQuery) ([]domain.RoomChange, error) {
	return s.repository.GetRoomChanges(query)
}

func (s *AcademicService) GetReservations(query domain.ReservationQuery) ([]domain.Reservation, error) {
	return s.repository.GetReservations(query)
}
