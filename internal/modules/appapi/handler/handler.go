// Package handler は App API（app.v1）の Connect ハンドラーを実装する。
// データの取得は app-bff-api と共通の service 層（internal/modules/app/service）に委譲する。
package handler

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// AnnouncementService はおしらせを取得する。
type AnnouncementService interface {
	GetAnnouncements() ([]domain.Announcement, error)
}

// AcademicService は学務情報を取得・更新する。
type AcademicService interface {
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

// UserService はユーザーと FCM トークンを取得・更新する。
type UserService interface {
	GetUser(id string) (*domain.User, error)
	UpsertUser(id string, req domain.UserRequest) (*domain.User, error)
	UpsertFCMToken(userID string, req domain.FCMTokenRequest) (*domain.FCMToken, error)
}
