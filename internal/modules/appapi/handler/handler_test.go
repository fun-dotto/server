package handler

import (
	"context"
	"testing"
	"time"
	"uuid"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func userCtx(id string) context.Context {
	return auth.WithUser(context.Background(), auth.User{ID: id})
}

func wantCode(t *testing.T, err error, code connect.Code) {
	t.Helper()
	if got := connect.CodeOf(err); got != code {
		t.Fatalf("code = %v, want %v (err: %v)", got, code, err)
	}
}

// DB に到達する前に返るエラーを検証する（db は nil）。
func TestValidationBeforeDB(t *testing.T) {
	h := New(nil, academiccalendar.Calendar{})
	ctx := userCtx("user1")
	subjectID := uuid.New().String()

	t.Run("他のユーザーの履修情報は削除できない", func(t *testing.T) {
		_, err := h.DeleteCourseRegistration(ctx, connect.NewRequest(&appv1.DeleteCourseRegistrationRequest{Id: "user2|" + subjectID}))
		wantCode(t, err, connect.CodePermissionDenied)
	})
	t.Run("不正な履修情報IDは PERMISSION_DENIED", func(t *testing.T) {
		_, err := h.DeleteCourseRegistration(ctx, connect.NewRequest(&appv1.DeleteCourseRegistrationRequest{Id: "broken"}))
		wantCode(t, err, connect.CodePermissionDenied)
	})
	t.Run("履修情報IDは必須", func(t *testing.T) {
		_, err := h.DeleteCourseRegistration(ctx, connect.NewRequest(&appv1.DeleteCourseRegistrationRequest{}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
	t.Run("未認証の場合は UNAUTHENTICATED", func(t *testing.T) {
		_, err := h.DeleteCourseRegistration(context.Background(), connect.NewRequest(&appv1.DeleteCourseRegistrationRequest{Id: "user1|" + subjectID}))
		wantCode(t, err, connect.CodeUnauthenticated)
	})
	t.Run("科目IDが不正", func(t *testing.T) {
		_, err := h.CreateCourseRegistration(ctx, connect.NewRequest(&appv1.CreateCourseRegistrationRequest{SubjectId: "x"}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
	t.Run("開講時期は必須", func(t *testing.T) {
		_, err := h.ListCourseRegistrations(ctx, connect.NewRequest(&appv1.ListCourseRegistrationsRequest{}))
		wantCode(t, err, connect.CodeInvalidArgument)
		_, err = h.ListTimetableItems(ctx, connect.NewRequest(&appv1.ListTimetableItemsRequest{
			Semesters: []appv1.SubjectSemester{appv1.SubjectSemester_SUBJECT_SEMESTER_UNSPECIFIED},
		}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
	t.Run("不正な日付", func(t *testing.T) {
		_, err := h.ListPersonalCalendarItems(ctx, connect.NewRequest(&appv1.ListPersonalCalendarItemsRequest{
			Dates: []*date.Date{{Year: 2026, Month: 2, Day: 30}},
		}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
	t.Run("FCM トークンは必須", func(t *testing.T) {
		_, err := h.RegisterFcmToken(ctx, connect.NewRequest(&appv1.RegisterFcmTokenRequest{}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
	t.Run("update_mask に未知のフィールド", func(t *testing.T) {
		_, err := h.UpdateUser(ctx, connect.NewRequest(&appv1.UpdateUserRequest{UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email"}}}))
		wantCode(t, err, connect.CodeInvalidArgument)
	})
}

func TestCourseRegistrationID(t *testing.T) {
	sid := uuid.New()
	id := encodeCourseRegistrationID("user1", sid)
	userID, got, ok := decodeCourseRegistrationID(id)
	if !ok || userID != "user1" || got != sid {
		t.Fatalf("decode(%q) = %q, %v, %v", id, userID, got, ok)
	}
}

func TestEscapeLike(t *testing.T) {
	if got := escapeLike(`a%b_c\`); got != `a\%b\_c\\` {
		t.Fatalf("escapeLike = %q", got)
	}
}

func strp(s string) *string { return &s }

func TestBuildPersonalCalendar(t *testing.T) {
	math := &model.Subject{Common: model.Common{ID: uuid.New()}, Name: "数学"}
	eng := &model.Subject{Common: model.Common{ID: uuid.New()}, Name: "英語"}
	other := &model.Subject{Common: model.Common{ID: uuid.New()}, Name: "未履修"}
	room1 := &model.Room{Common: model.Common{ID: uuid.New()}, Name: "101"}
	room2 := &model.Room{Common: model.Common{ID: uuid.New()}, Name: "202"}

	item := func(s *model.Subject, dow, period string, rooms ...*model.Room) model.TimetableItem {
		ti := model.TimetableItem{SubjectID: s.ID, Subject: s, DayOfWeek: strp(dow), Period: strp(period)}
		for _, r := range rooms {
			ti.Rooms = append(ti.Rooms, model.TimetableItemRoom{Room: r})
		}
		return ti
	}
	d := func(day int) time.Time { return time.Date(2026, 4, day, 0, 0, 0, 0, time.UTC) }

	// 2026-04-13 月曜、2026-04-14 火曜（月曜振替）、2026-04-15 水曜（休日）
	cal := academiccalendar.Calendar{
		SubstituteDays: map[string]time.Weekday{"2026-04-14": time.Monday},
		Holidays:       map[string]struct{}{"2026-04-15": {}},
	}
	in := calendarInput{
		registeredSubjectIDs: map[string]struct{}{math.ID.String(): {}, eng.ID.String(): {}},
		timetableItems: []model.TimetableItem{
			item(math, "Monday", "Period1", room1),
			item(eng, "Monday", "Period2", room1),
			item(other, "Monday", "Period3", room1),
			item(math, "Wednesday", "Period1", room2),
		},
		cancelledClasses: []model.CancelledClass{{SubjectID: math.ID, Subject: math, Date: d(13), Period: "Period1"}},
		roomChanges:      []model.RoomChange{{SubjectID: eng.ID, Subject: eng, Date: d(13), Period: "Period2", NewRoom: room2}},
		makeupClasses: []model.MakeupClass{
			{SubjectID: math.ID, Subject: math, Date: d(14), Period: "Period5"},
			// 休講と同じ枠の補講は休講を優先する
			{SubjectID: math.ID, Subject: math, Date: d(13), Period: "Period1"},
		},
	}

	got := buildPersonalCalendar(cal, []time.Time{d(13), d(14), d(15)}, in)

	type row struct {
		day     int
		period  string
		subject string
		state   appv1.PersonalCalendarItem_State
		room    string
	}
	want := []row{
		{13, "Period1", "数学", appv1.PersonalCalendarItem_STATE_CANCELLED, "101"},
		{13, "Period2", "英語", appv1.PersonalCalendarItem_STATE_ROOM_CHANGED, "202"},
		{14, "Period1", "数学", appv1.PersonalCalendarItem_STATE_NORMAL, "101"},
		{14, "Period2", "英語", appv1.PersonalCalendarItem_STATE_NORMAL, "101"},
		{14, "Period5", "数学", appv1.PersonalCalendarItem_STATE_MAKEUP, "101"},
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		g := got[i]
		room := ""
		if len(g.rooms) > 0 {
			room = g.rooms[0].Name
		}
		if g.date.Day() != w.day || g.period != w.period || g.subject.Name != w.subject || g.state != w.state || room != w.room {
			t.Errorf("[%d] = {%d %s %s %v %s}, want %+v", i, g.date.Day(), g.period, g.subject.Name, g.state, room, w)
		}
	}
}

func TestBuildPersonalCalendarWithoutRegistrations(t *testing.T) {
	got := buildPersonalCalendar(academiccalendar.Calendar{}, []time.Time{time.Now()}, calendarInput{})
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestTermsOf(t *testing.T) {
	terms := []academiccalendar.Term{
		{Year: 2026, Semester: "Q1", Start: "2026-04-06", End: "2026-06-04"},
		{Year: 2026, Semester: "H1", Start: "2026-04-06", End: "2026-08-03"},
	}
	year, semesters := termsOf(terms, []time.Time{time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)})
	if year == nil || *year != 2026 || len(semesters) != 2 || semesters[0] != "H1" || semesters[1] != "Q1" {
		t.Fatalf("termsOf = %v, %v", year, semesters)
	}
}
