package service

import (
	"context"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	userdomain "github.com/fun-dotto/server/internal/modules/user/domain"
)

type CancelledClassRepository interface {
	List(ctx context.Context, filter academicdomain.CancelledClassListFilter) ([]academicdomain.CancelledClass, error)
}

type MakeupClassRepository interface {
	List(ctx context.Context, filter academicdomain.MakeupClassListFilter) ([]academicdomain.MakeupClass, error)
}

type RoomChangeRepository interface {
	List(ctx context.Context, filter academicdomain.RoomChangeListFilter) ([]academicdomain.RoomChange, error)
}

type CourseRegistrationRepository interface {
	ListUserIDsBySubject(ctx context.Context, subjectID string) ([]string, error)
}

type NotificationRepository interface {
	UpsertNotification(ctx context.Context, notification userdomain.Notification) (userdomain.Notification, error)
}

type ClassChangeNotificationService struct {
	cancelled    CancelledClassRepository
	makeup       MakeupClassRepository
	roomChange   RoomChangeRepository
	courseReg    CourseRegistrationRepository
	notification NotificationRepository
}

func NewClassChangeNotificationService(
	cancelled CancelledClassRepository,
	makeup MakeupClassRepository,
	roomChange RoomChangeRepository,
	courseReg CourseRegistrationRepository,
	notification NotificationRepository,
) *ClassChangeNotificationService {
	return &ClassChangeNotificationService{
		cancelled:    cancelled,
		makeup:       makeup,
		roomChange:   roomChange,
		courseReg:    courseReg,
		notification: notification,
	}
}
