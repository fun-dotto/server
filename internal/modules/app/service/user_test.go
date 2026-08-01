package service

import (
	"testing"
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name     string
		repo     UserRepository
		userID   string
		validate func(t *testing.T, user *domain.User, err error)
	}{
		{
			name:   "正常系: ユーザーを取得できる",
			repo:   repository.NewMockUserRepository(),
			userID: "user1",
			validate: func(t *testing.T, user *domain.User, err error) {
				require.NoError(t, err)
				assert.Equal(t, "user1", user.ID)
				assert.Equal(t, "user1@example.com", user.Email)
				assert.NotNil(t, user.Grade)
				assert.NotNil(t, user.Course)
				assert.NotNil(t, user.Class)
			},
		},
		{
			name:   "異常系: 存在しないユーザーIDの場合エラーを返す",
			repo:   repository.NewMockUserRepository(),
			userID: "nonexistent",
			validate: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, domain.ErrUserNotFound)
			},
		},
		{
			name:   "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo:   repository.NewMockUserRepositoryWithError("getUser", assert.AnError),
			userID: "user1",
			validate: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUserService(tt.repo)
			user, err := svc.GetUser(tt.userID)
			tt.validate(t, user, err)
		})
	}
}

func TestUserService_UpsertUser(t *testing.T) {
	grade := domain.GradeB1
	course := domain.CourseInformationDesign
	class := domain.ClassB

	tests := []struct {
		name     string
		repo     UserRepository
		userID   string
		req      domain.UserRequest
		validate func(t *testing.T, user *domain.User, err error)
	}{
		{
			name:   "正常系: 新規ユーザーを作成できる",
			repo:   repository.NewMockUserRepositoryEmpty(),
			userID: "newuser",
			req: domain.UserRequest{
				Email:  "new@example.com",
				Grade:  &grade,
				Course: &course,
				Class:  &class,
			},
			validate: func(t *testing.T, user *domain.User, err error) {
				require.NoError(t, err)
				assert.Equal(t, "newuser", user.ID)
				assert.Equal(t, "new@example.com", user.Email)
				require.NotNil(t, user.Grade)
				assert.Equal(t, grade, *user.Grade)
				require.NotNil(t, user.Course)
				assert.Equal(t, course, *user.Course)
				require.NotNil(t, user.Class)
				assert.Equal(t, class, *user.Class)
			},
		},
		{
			name:   "正常系: 既存ユーザーを更新できる",
			repo:   repository.NewMockUserRepository(),
			userID: "user1",
			req: domain.UserRequest{
				Email:  "updated@example.com",
				Grade:  &grade,
				Course: &course,
				Class:  &class,
			},
			validate: func(t *testing.T, user *domain.User, err error) {
				require.NoError(t, err)
				assert.Equal(t, "user1", user.ID)
				assert.Equal(t, "updated@example.com", user.Email)
			},
		},
		{
			name:   "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo:   repository.NewMockUserRepositoryWithError("upsert", assert.AnError),
			userID: "user1",
			req: domain.UserRequest{
				Email: "test@example.com",
			},
			validate: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUserService(tt.repo)
			user, err := svc.UpsertUser(tt.userID, tt.req)
			tt.validate(t, user, err)
		})
	}
}

func TestUserService_UpsertFCMToken(t *testing.T) {
	tests := []struct {
		name     string
		repo     UserRepository
		userID   string
		req      domain.FCMTokenRequest
		validate func(t *testing.T, token *domain.FCMToken, err error)
	}{
		{
			name:   "正常系: FCMトークンを作成・更新できる",
			repo:   repository.NewMockUserRepository(),
			userID: "user1",
			req:    domain.FCMTokenRequest{Token: "fcm-token-1"},
			validate: func(t *testing.T, token *domain.FCMToken, err error) {
				require.NoError(t, err)
				assert.Equal(t, "fcm-token-1", token.Token)
				assert.WithinDuration(t, time.Now(), token.CreatedAt, time.Minute)
				assert.WithinDuration(t, time.Now(), token.UpdatedAt, time.Minute)
			},
		},
		{
			name:   "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo:   repository.NewMockUserRepositoryWithError("upsertFCMToken", assert.AnError),
			userID: "user1",
			req:    domain.FCMTokenRequest{Token: "fcm-token-1"},
			validate: func(t *testing.T, token *domain.FCMToken, err error) {
				require.Error(t, err)
				assert.Nil(t, token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUserService(tt.repo)
			token, err := svc.UpsertFCMToken(tt.userID, tt.req)
			tt.validate(t, token, err)
		})
	}
}
