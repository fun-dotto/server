package service

import (
	"testing"

	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcademicService_GetFaculties(t *testing.T) {
	tests := []struct {
		name     string
		repo     AcademicRepository
		validate func(t *testing.T, faculties []domain.Faculty, err error)
	}{
		{
			name: "正常系: 教員一覧を取得できる",
			repo: repository.NewMockAcademicRepository(),
			validate: func(t *testing.T, faculties []domain.Faculty, err error) {
				require.NoError(t, err)
				assert.Len(t, faculties, 1)
				assert.Equal(t, "f1", faculties[0].ID)
				assert.Equal(t, "田中太郎", faculties[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			faculties, err := svc.GetFaculties()
			tt.validate(t, faculties, err)
		})
	}
}

func TestAcademicService_GetFaculty(t *testing.T) {
	tests := []struct {
		name     string
		repo     AcademicRepository
		id       string
		validate func(t *testing.T, faculty *domain.Faculty, err error)
	}{
		{
			name: "正常系: 教員を取得できる",
			repo: repository.NewMockAcademicRepository(),
			id:   "f1",
			validate: func(t *testing.T, faculty *domain.Faculty, err error) {
				require.NoError(t, err)
				assert.Equal(t, "f1", faculty.ID)
				assert.Equal(t, "田中太郎", faculty.Name)
			},
		},
		{
			name: "異常系: 存在しないIDの場合エラーを返す",
			repo: repository.NewMockAcademicRepository(),
			id:   "nonexistent",
			validate: func(t *testing.T, faculty *domain.Faculty, err error) {
				require.Error(t, err)
				assert.Nil(t, faculty)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			faculty, err := svc.GetFaculty(tt.id)
			tt.validate(t, faculty, err)
		})
	}
}

func TestAcademicService_GetSubjects(t *testing.T) {
	tests := []struct {
		name     string
		repo     AcademicRepository
		query    domain.SubjectQuery
		validate func(t *testing.T, subjects []domain.Subject, err error)
	}{
		{
			name:  "正常系: 科目一覧を取得しFaculty情報が補完される",
			repo:  repository.NewMockAcademicRepository(),
			query: domain.SubjectQuery{},
			validate: func(t *testing.T, subjects []domain.Subject, err error) {
				require.NoError(t, err)
				assert.Len(t, subjects, 1)
				assert.Equal(t, "s1", subjects[0].ID)
				assert.Equal(t, "プログラミング基礎", subjects[0].Name)
				require.Len(t, subjects[0].Faculties, 1)
				assert.Equal(t, "田中太郎", subjects[0].Faculties[0].Faculty.Name)
				assert.Equal(t, "tanaka@example.com", subjects[0].Faculties[0].Faculty.Email)
			},
		},
		{
			name:  "異常系: GetSubjectsがエラーを返す場合エラーを返す",
			repo:  repository.NewMockAcademicRepositoryWithError("getSubjects", assert.AnError),
			query: domain.SubjectQuery{},
			validate: func(t *testing.T, subjects []domain.Subject, err error) {
				require.Error(t, err)
				assert.Nil(t, subjects)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			subjects, err := svc.GetSubjects(tt.query)
			tt.validate(t, subjects, err)
		})
	}
}

func TestAcademicService_GetSubject(t *testing.T) {
	tests := []struct {
		name     string
		repo     AcademicRepository
		id       string
		validate func(t *testing.T, subject *domain.Subject, err error)
	}{
		{
			name: "正常系: 科目を取得しFaculty情報が補完される",
			repo: repository.NewMockAcademicRepository(),
			id:   "s1",
			validate: func(t *testing.T, subject *domain.Subject, err error) {
				require.NoError(t, err)
				assert.Equal(t, "s1", subject.ID)
				assert.Equal(t, "プログラミング基礎", subject.Name)
				require.Len(t, subject.Faculties, 1)
				assert.Equal(t, "田中太郎", subject.Faculties[0].Faculty.Name)
				assert.Equal(t, "tanaka@example.com", subject.Faculties[0].Faculty.Email)
			},
		},
		{
			name: "異常系: 存在しないIDの場合エラーを返す",
			repo: repository.NewMockAcademicRepository(),
			id:   "nonexistent",
			validate: func(t *testing.T, subject *domain.Subject, err error) {
				require.Error(t, err)
				assert.Nil(t, subject)
			},
		},
		{
			name: "異常系: GetSubjectがエラーを返す場合エラーを返す",
			repo: repository.NewMockAcademicRepositoryWithError("getSubject", assert.AnError),
			id:   "s1",
			validate: func(t *testing.T, subject *domain.Subject, err error) {
				require.Error(t, err)
				assert.Nil(t, subject)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			subject, err := svc.GetSubject(tt.id)
			tt.validate(t, subject, err)
		})
	}
}

func TestAcademicService_GetCourseRegistrations(t *testing.T) {
	year := 2026

	tests := []struct {
		name      string
		repo      AcademicRepository
		userID    string
		semesters []domain.CourseSemester
		year      *int
		validate  func(t *testing.T, registrations []domain.CourseRegistration, err error)
	}{
		{
			name:      "正常系: 履修登録一覧を取得できる",
			repo:      repository.NewMockAcademicRepository(),
			userID:    "user1",
			semesters: []domain.CourseSemester{domain.CourseSemesterQ1},
			year:      &year,
			validate: func(t *testing.T, registrations []domain.CourseRegistration, err error) {
				require.NoError(t, err)
				assert.Len(t, registrations, 1)
				assert.Equal(t, "cr1", registrations[0].ID)
			},
		},
		{
			name:      "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo:      repository.NewMockAcademicRepositoryWithError("getRegistrations", assert.AnError),
			userID:    "user1",
			semesters: []domain.CourseSemester{domain.CourseSemesterQ1},
			year:      &year,
			validate: func(t *testing.T, registrations []domain.CourseRegistration, err error) {
				require.Error(t, err)
				assert.Nil(t, registrations)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			registrations, err := svc.GetCourseRegistrations(tt.userID, tt.semesters, tt.year)
			tt.validate(t, registrations, err)
		})
	}
}

func TestAcademicService_CreateCourseRegistration(t *testing.T) {
	tests := []struct {
		name      string
		repo      AcademicRepository
		userID    string
		subjectID string
		validate  func(t *testing.T, registration *domain.CourseRegistration, err error)
	}{
		{
			name:      "正常系: 履修登録を作成できる",
			repo:      repository.NewMockAcademicRepository(),
			userID:    "user1",
			subjectID: "s1",
			validate: func(t *testing.T, registration *domain.CourseRegistration, err error) {
				require.NoError(t, err)
				assert.Equal(t, "cr-new", registration.ID)
			},
		},
		{
			name:      "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo:      repository.NewMockAcademicRepositoryWithError("create", assert.AnError),
			userID:    "user1",
			subjectID: "s1",
			validate: func(t *testing.T, registration *domain.CourseRegistration, err error) {
				require.Error(t, err)
				assert.Nil(t, registration)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			registration, err := svc.CreateCourseRegistration(tt.userID, tt.subjectID)
			tt.validate(t, registration, err)
		})
	}
}

func TestAcademicService_DeleteCourseRegistration(t *testing.T) {
	tests := []struct {
		name     string
		repo     AcademicRepository
		id       string
		validate func(t *testing.T, err error)
	}{
		{
			name: "正常系: 履修登録を削除できる",
			repo: repository.NewMockAcademicRepository(),
			id:   "cr1",
			validate: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAcademicRepositoryWithError("delete", assert.AnError),
			id:   "cr1",
			validate: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAcademicService(tt.repo)
			err := svc.DeleteCourseRegistration(tt.id)
			tt.validate(t, err)
		})
	}
}
