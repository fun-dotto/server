package handler

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func userCtx() context.Context {
	return auth.WithUser(context.Background(), auth.User{ID: "user1", Email: "user1@example.com"})
}

func academic() *service.AcademicService {
	return service.NewAcademicService(repository.NewMockAcademicRepository())
}

func wantCode(t *testing.T, err error, code connect.Code) {
	t.Helper()
	if got := connect.CodeOf(err); got != code {
		t.Fatalf("code = %v, want %v (err: %v)", got, code, err)
	}
}

func TestListSubjects(t *testing.T) {
	h := NewSubjectHandler(academic())
	res, err := h.ListSubjects(context.Background(), connect.NewRequest(&appv1.ListSubjectsRequest{
		Grades: []appv1.Grade{appv1.Grade_GRADE_B1},
	}))
	if err != nil {
		t.Fatal(err)
	}
	s := res.Msg.GetSubjects()
	if len(s) != 1 || s[0].GetId() != "s1" || s[0].GetSemester() != appv1.SubjectSemester_SUBJECT_SEMESTER_Q1 {
		t.Fatalf("unexpected subjects: %v", s)
	}
	if s[0].GetSyllabus() != nil || len(s[0].GetRequirements()) != 0 {
		t.Fatal("list must not set detail fields")
	}

	_, err = h.ListSubjects(context.Background(), connect.NewRequest(&appv1.ListSubjectsRequest{
		Grades: []appv1.Grade{appv1.Grade_GRADE_UNSPECIFIED},
	}))
	wantCode(t, err, connect.CodeInvalidArgument)
}

func TestGetSubject(t *testing.T) {
	h := NewSubjectHandler(academic())
	res, err := h.GetSubject(context.Background(), connect.NewRequest(&appv1.GetSubjectRequest{Id: "s1"}))
	if err != nil {
		t.Fatal(err)
	}
	if got := res.Msg.GetSubject().GetRequirements(); len(got) != 1 || got[0].GetCourse() != appv1.Course_COURSE_INFORMATION_SYSTEM {
		t.Fatalf("unexpected requirements: %v", got)
	}

	_, err = h.GetSubject(context.Background(), connect.NewRequest(&appv1.GetSubjectRequest{}))
	wantCode(t, err, connect.CodeInvalidArgument)

	_, err = NewSubjectHandler(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getSubject", errors.New("boom")))).
		GetSubject(context.Background(), connect.NewRequest(&appv1.GetSubjectRequest{Id: "s1"}))
	wantCode(t, err, connect.CodeInternal)
}

func TestListCourseRegistrations(t *testing.T) {
	h := NewCourseRegistrationHandler(academic())

	_, err := h.ListCourseRegistrations(context.Background(), connect.NewRequest(&appv1.ListCourseRegistrationsRequest{
		Semesters: []appv1.SubjectSemester{appv1.SubjectSemester_SUBJECT_SEMESTER_Q1},
	}))
	wantCode(t, err, connect.CodeUnauthenticated)

	_, err = h.ListCourseRegistrations(userCtx(), connect.NewRequest(&appv1.ListCourseRegistrationsRequest{}))
	wantCode(t, err, connect.CodeInvalidArgument)

	res, err := h.ListCourseRegistrations(userCtx(), connect.NewRequest(&appv1.ListCourseRegistrationsRequest{
		Semesters: []appv1.SubjectSemester{appv1.SubjectSemester_SUBJECT_SEMESTER_Q1},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Msg.GetCourseRegistrations()) != 1 {
		t.Fatalf("unexpected registrations: %v", res.Msg.GetCourseRegistrations())
	}
}

func TestListPersonalCalendarItems(t *testing.T) {
	h := NewPersonalCalendarItemHandler(academic())

	_, err := h.ListPersonalCalendarItems(userCtx(), connect.NewRequest(&appv1.ListPersonalCalendarItemsRequest{
		Dates: []*date.Date{{Year: 2026, Month: 2, Day: 30}},
	}))
	wantCode(t, err, connect.CodeInvalidArgument)

	res, err := h.ListPersonalCalendarItems(userCtx(), connect.NewRequest(&appv1.ListPersonalCalendarItemsRequest{
		Dates: []*date.Date{{Year: 2026, Month: 4, Day: 1}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	items := res.Msg.GetPersonalCalendarItems()
	if len(items) != 1 || items[0].GetState() != appv1.PersonalCalendarItem_STATE_NORMAL || items[0].GetDate().GetDay() != 1 {
		t.Fatalf("unexpected items: %v", items)
	}
}

func TestUser(t *testing.T) {
	repo := repository.NewMockUserRepository()
	h := NewUserHandler(service.NewUserService(repo))

	res, err := h.GetUser(userCtx(), connect.NewRequest(&appv1.GetUserRequest{}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Msg.GetUser().GetGrade() != appv1.Grade_GRADE_B2 {
		t.Fatalf("unexpected user: %v", res.Msg.GetUser())
	}

	// update_mask で指定したフィールドのみ更新し、他は現在の値を引き継ぐ。
	grade := appv1.Grade_GRADE_B3
	upd, err := h.UpdateUser(userCtx(), connect.NewRequest(&appv1.UpdateUserRequest{
		Grade:      &grade,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"grade"}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	u := upd.Msg.GetUser()
	if u.GetGrade() != appv1.Grade_GRADE_B3 || u.GetCourse() != appv1.Course_COURSE_INFORMATION_SYSTEM || u.GetClass() != appv1.Class_CLASS_A {
		t.Fatalf("unexpected user: %v", u)
	}

	_, err = h.UpdateUser(userCtx(), connect.NewRequest(&appv1.UpdateUserRequest{
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email"}},
	}))
	wantCode(t, err, connect.CodeInvalidArgument)

	_, err = NewUserHandler(service.NewUserService(repository.NewMockUserRepositoryEmpty())).
		GetUser(userCtx(), connect.NewRequest(&appv1.GetUserRequest{}))
	wantCode(t, err, connect.CodeNotFound)

	// ユーザーが存在しない場合は作成する。
	created, err := NewUserHandler(service.NewUserService(repository.NewMockUserRepositoryEmpty())).
		UpdateUser(userCtx(), connect.NewRequest(&appv1.UpdateUserRequest{Grade: &grade}))
	if err != nil {
		t.Fatal(err)
	}
	if created.Msg.GetUser().GetGrade() != appv1.Grade_GRADE_B3 {
		t.Fatalf("unexpected user: %v", created.Msg.GetUser())
	}
}

func TestRegisterFcmToken(t *testing.T) {
	h := NewFcmTokenHandler(service.NewUserService(repository.NewMockUserRepository()))
	_, err := h.RegisterFcmToken(userCtx(), connect.NewRequest(&appv1.RegisterFcmTokenRequest{}))
	wantCode(t, err, connect.CodeInvalidArgument)

	res, err := h.RegisterFcmToken(userCtx(), connect.NewRequest(&appv1.RegisterFcmTokenRequest{Token: "tok"}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Msg.GetFcmToken().GetToken() != "tok" {
		t.Fatalf("unexpected token: %v", res.Msg.GetFcmToken())
	}
}

type fakeAnnouncements []domain.Announcement

func (f fakeAnnouncements) GetAnnouncements() ([]domain.Announcement, error) { return f, nil }

func TestPagination(t *testing.T) {
	items := make(fakeAnnouncements, 5)
	for i := range items {
		items[i] = domain.Announcement{ID: string(rune('a' + i))}
	}
	h := NewAnnouncementHandler(items)

	var ids []string
	token := ""
	for range 10 {
		res, err := h.ListAnnouncements(context.Background(), connect.NewRequest(&appv1.ListAnnouncementsRequest{PageSize: 2, PageToken: token}))
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range res.Msg.GetAnnouncements() {
			ids = append(ids, a.GetId())
		}
		token = res.Msg.GetNextPageToken()
		if token == "" {
			break
		}
	}
	if got := len(ids); got != 5 || ids[0] != "a" || ids[4] != "e" {
		t.Fatalf("ids = %v", ids)
	}

	_, err := h.ListAnnouncements(context.Background(), connect.NewRequest(&appv1.ListAnnouncementsRequest{PageSize: -1}))
	wantCode(t, err, connect.CodeInvalidArgument)
	_, err = h.ListAnnouncements(context.Background(), connect.NewRequest(&appv1.ListAnnouncementsRequest{PageToken: "broken"}))
	wantCode(t, err, connect.CodeInvalidArgument)
}

func TestPaginationRejectsChangedFilter(t *testing.T) {
	h := NewSubjectHandler(service.NewAcademicService(&manySubjects{}))
	res, err := h.ListSubjects(context.Background(), connect.NewRequest(&appv1.ListSubjectsRequest{PageSize: 1, Query: "a"}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Msg.GetNextPageToken() == "" {
		t.Fatal("expected next page token")
	}
	_, err = h.ListSubjects(context.Background(), connect.NewRequest(&appv1.ListSubjectsRequest{PageSize: 1, Query: "b", PageToken: res.Msg.GetNextPageToken()}))
	wantCode(t, err, connect.CodeInvalidArgument)
}

type manySubjects struct{ *repository.MockAcademicRepository }

func (*manySubjects) GetSubjects(domain.SubjectQuery) ([]domain.Subject, error) {
	return []domain.Subject{{ID: "1"}, {ID: "2"}}, nil
}
