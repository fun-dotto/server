// Package handler は App API（app.v1）の Connect ハンドラーを実装する。
// 他の API には依存せず、DB（internal/shared/model）を直接操作する。
package handler

import (
	"context"
	"time"

	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// Handler は app.v1 の各サービスを実装する。
type Handler struct {
	db       *gorm.DB
	calendar academiccalendar.Calendar
	now      func() time.Time
}

var (
	_ appv1connect.AnnouncementServiceHandler         = (*Handler)(nil)
	_ appv1connect.CancelledClassServiceHandler       = (*Handler)(nil)
	_ appv1connect.CourseRegistrationServiceHandler   = (*Handler)(nil)
	_ appv1connect.FcmTokenServiceHandler             = (*Handler)(nil)
	_ appv1connect.MakeupClassServiceHandler          = (*Handler)(nil)
	_ appv1connect.PersonalCalendarItemServiceHandler = (*Handler)(nil)
	_ appv1connect.ReservationServiceHandler          = (*Handler)(nil)
	_ appv1connect.RoomChangeServiceHandler           = (*Handler)(nil)
	_ appv1connect.SubjectServiceHandler              = (*Handler)(nil)
	_ appv1connect.TimetableItemServiceHandler        = (*Handler)(nil)
	_ appv1connect.UserServiceHandler                 = (*Handler)(nil)
)

// New は Handler を返す。
func New(db *gorm.DB, calendar academiccalendar.Calendar) *Handler {
	return &Handler{db: db, calendar: calendar, now: time.Now}
}

// currentUserID は認証済みのユーザー ID を返す。
func currentUserID(ctx context.Context) (string, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return "", unauthenticated()
	}
	return user.ID, nil
}

// defaultYear は年度の指定がない場合に今年度を返す。
func (h *Handler) defaultYear(y *int32) int {
	if y != nil {
		return int(*y)
	}
	return academiccalendar.CurrentAcademicYear(h.now())
}

// attachRoomFaculties は今年度の教員室の教員を Room に設定する。
func (h *Handler) attachRoomFaculties(ctx context.Context, rooms ...*appv1.Room) error {
	ids := make([]string, 0, len(rooms))
	for _, r := range rooms {
		if r != nil {
			ids = append(ids, r.GetId())
		}
	}
	if len(ids) == 0 {
		return nil
	}
	var records []model.FacultyRoom
	if err := h.db.WithContext(ctx).
		Preload("Faculty").
		Where("room_id IN ? AND year = ?", ids, academiccalendar.CurrentAcademicYear(h.now())).
		Order("faculty_id ASC").
		Find(&records).Error; err != nil {
		return err
	}
	faculties := make(map[string]*appv1.Faculty, len(records))
	for _, rec := range records {
		key := rec.RoomID.String()
		if _, ok := faculties[key]; !ok && rec.Faculty != nil {
			faculties[key] = toProtoFaculty(rec.Faculty)
		}
	}
	for _, r := range rooms {
		if r == nil {
			continue
		}
		if f, ok := faculties[r.GetId()]; ok {
			r.Faculty = f
		}
	}
	return nil
}

func preload(db *gorm.DB, names ...string) *gorm.DB {
	for _, n := range names {
		db = db.Preload(n)
	}
	return db
}
