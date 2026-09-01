package handler

import (
	"context"
	"errors"
	"testing"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

// fakeSubjectService は List に渡された filter を記録するだけのスタブ。
type fakeSubjectService struct {
	gotFilter domain.SubjectListFilter
	called    bool
}

func (f *fakeSubjectService) List(_ context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error) {
	f.gotFilter = filter
	f.called = true
	return nil, nil
}

func (f *fakeSubjectService) GetByID(context.Context, string) (domain.Subject, error) {
	return domain.Subject{}, nil
}

func (f *fakeSubjectService) Delete(context.Context, string) error { return nil }

func (f *fakeSubjectService) GetSyllabus(context.Context, string) (domain.Syllabus, error) {
	return domain.Syllabus{}, nil
}

type fakeUserService struct {
	user   domain.User
	found  bool
	err    error
	gotID  string
	called bool
}

func (f *fakeUserService) FindByID(_ context.Context, id string) (domain.User, bool, error) {
	f.gotID = id
	f.called = true
	return f.user, f.found, f.err
}

func TestSubjectsV1List_UserAttributeSort(t *testing.T) {
	course := domain.CourseTypeComplexSystem
	grade := domain.GradeB3
	userID := "firebase-uid-001"

	tests := []struct {
		name string
		// 入力
		paramUserID *string
		userSvc     *fakeUserService
		// 期待値
		wantUserSvcCalled bool
		wantSortEnabled   bool
		wantSortCourse    *domain.CourseType
		wantSortGrade     *domain.Grade
		wantErr           bool
	}{
		{
			name:              "userId 未指定ならソートせず UserService も呼ばない",
			paramUserID:       nil,
			userSvc:           &fakeUserService{},
			wantUserSvcCalled: false,
			wantSortEnabled:   false,
		},
		{
			name:        "userId 指定かつユーザーが存在すればソート条件が立つ",
			paramUserID: &userID,
			userSvc: &fakeUserService{
				user:  domain.User{ID: userID, Course: &course, Grade: &grade},
				found: true,
			},
			wantUserSvcCalled: true,
			wantSortEnabled:   true,
			wantSortCourse:    &course,
			wantSortGrade:     &grade,
		},
		{
			name:              "ユーザーが存在しなければエラーにせず現行順にフォールバックする",
			paramUserID:       &userID,
			userSvc:           &fakeUserService{found: false},
			wantUserSvcCalled: true,
			wantSortEnabled:   false,
		},
		{
			name:        "コース・学年が未設定のユーザーでもソートは有効になる",
			paramUserID: &userID,
			userSvc: &fakeUserService{
				user:  domain.User{ID: userID},
				found: true,
			},
			wantUserSvcCalled: true,
			wantSortEnabled:   true,
			wantSortCourse:    nil,
			wantSortGrade:     nil,
		},
		{
			name:              "UserService がエラーを返せばハンドラもエラーを返す",
			paramUserID:       &userID,
			userSvc:           &fakeUserService{err: errors.New("db is down")},
			wantUserSvcCalled: true,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subjectSvc := &fakeSubjectService{}
			h := &Handler{subjectSvc: subjectSvc, userSvc: tt.userSvc}

			_, err := h.SubjectsV1List(context.Background(), api.SubjectsV1ListRequestObject{
				Params: api.SubjectsV1ListParams{UserId: tt.paramUserID},
			})

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーを期待したが nil だった")
				}
				if subjectSvc.called {
					t.Error("エラー時は SubjectService.List を呼ぶべきではない")
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if tt.userSvc.called != tt.wantUserSvcCalled {
				t.Errorf("UserService.FindByID called = %v, want %v", tt.userSvc.called, tt.wantUserSvcCalled)
			}
			if tt.wantUserSvcCalled && tt.userSvc.gotID != userID {
				t.Errorf("FindByID に渡された id = %q, want %q", tt.userSvc.gotID, userID)
			}

			got := subjectSvc.gotFilter
			if got.SortByUserAttribute != tt.wantSortEnabled {
				t.Errorf("SortByUserAttribute = %v, want %v", got.SortByUserAttribute, tt.wantSortEnabled)
			}
			if !equalCourse(got.SortCourse, tt.wantSortCourse) {
				t.Errorf("SortCourse = %v, want %v", got.SortCourse, tt.wantSortCourse)
			}
			if !equalGrade(got.SortGrade, tt.wantSortGrade) {
				t.Errorf("SortGrade = %v, want %v", got.SortGrade, tt.wantSortGrade)
			}
		})
	}
}

func equalCourse(a, b *domain.CourseType) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func equalGrade(a, b *domain.Grade) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
