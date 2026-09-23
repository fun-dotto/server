package handler

import (
	"context"
	"fmt"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubjectsV1List(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		validate func(t *testing.T, resp api.SubjectsV1ListResponseObject, err error)
	}{
		{
			name:    "正常に科目一覧が取得できる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.SubjectsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.SubjectsV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.Len(t, result.Subjects, 1)
				assert.Equal(t, "s1", result.Subjects[0].Id)
				assert.Equal(t, "プログラミング基礎", result.Subjects[0].Name)
			},
		},
		{
			name:    "科目のFacultyが正しくマッピングされる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			validate: func(t *testing.T, resp api.SubjectsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.SubjectsV1List200JSONResponse)
				require.True(t, ok)
				require.Len(t, result.Subjects, 1)
				require.Len(t, result.Subjects[0].Faculties, 1)
				assert.Equal(t, "f1", result.Subjects[0].Faculties[0].Faculty.Id)
				assert.Equal(t, "田中太郎", result.Subjects[0].Faculties[0].Faculty.Name)
				assert.True(t, result.Subjects[0].Faculties[0].IsPrimary)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			validate: func(t *testing.T, resp api.SubjectsV1ListResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getSubjects", fmt.Errorf("db error"))))),
			validate: func(t *testing.T, resp api.SubjectsV1ListResponseObject, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.SubjectsV1List(context.Background(), api.SubjectsV1ListRequestObject{})
			tt.validate(t, resp, err)
		})
	}
}

func TestSubjectsV1Detail(t *testing.T) {
	tests := []struct {
		name     string
		handler  *Handler
		request  api.SubjectsV1DetailRequestObject
		validate func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error)
	}{
		{
			name:    "正常に科目詳細が取得できる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			request: api.SubjectsV1DetailRequestObject{Id: "s1"},
			validate: func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.SubjectsV1Detail200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.Equal(t, "s1", result.Subject.Id)
				assert.Equal(t, "プログラミング基礎", result.Subject.Name)
				assert.Equal(t, 2026, result.Subject.Year)
				assert.Equal(t, 2, result.Subject.Credit)
				assert.Equal(t, api.DottoFoundationV1CourseSemester("Q1"), result.Subject.Semester)
			},
		},
		{
			name:    "科目詳細のRequirementsが正しくマッピングされる",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			request: api.SubjectsV1DetailRequestObject{Id: "s1"},
			validate: func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.SubjectsV1Detail200JSONResponse)
				require.True(t, ok)
				require.Len(t, result.Subject.Requirements, 1)
				assert.Equal(t, api.DottoFoundationV1Course("InformationSystem"), result.Subject.Requirements[0].Course)
				assert.Equal(t, api.DottoFoundationV1SubjectRequirementType("Required"), result.Subject.Requirements[0].RequirementType)
			},
		},
		{
			name:    "academicServiceが未設定の場合エラーを返す",
			handler: NewHandler(),
			request: api.SubjectsV1DetailRequestObject{Id: "s1"},
			validate: func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errAcademicServiceNotConfigured)
			},
		},
		{
			name:    "存在しない科目IDの場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepository()))),
			request: api.SubjectsV1DetailRequestObject{Id: "nonexistent"},
			validate: func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get subject")
			},
		},
		{
			name:    "serviceがエラーを返す場合エラーを返す",
			handler: NewHandler(WithAcademicService(service.NewAcademicService(repository.NewMockAcademicRepositoryWithError("getSubject", fmt.Errorf("db error"))))),
			request: api.SubjectsV1DetailRequestObject{Id: "s1"},
			validate: func(t *testing.T, resp api.SubjectsV1DetailResponseObject, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler.SubjectsV1Detail(context.Background(), tt.request)
			tt.validate(t, resp, err)
		})
	}
}
