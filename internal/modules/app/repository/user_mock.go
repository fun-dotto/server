package repository

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// MockUserRepository はテスト用のモック
type MockUserRepository struct {
	users             map[string]domain.User
	upsertErr         error
	getUserErr        error
	upsertFCMTokenErr error
	fcmTokens         map[string]domain.FCMToken
}

func NewMockUserRepository() *MockUserRepository {
	grade := domain.GradeB2
	course := domain.CourseInformationSystem
	class := domain.ClassA

	return &MockUserRepository{
		users: map[string]domain.User{
			"user1": {
				ID:     "user1",
				Email:  "user1@example.com",
				Grade:  &grade,
				Course: &course,
				Class:  &class,
			},
		},
		fcmTokens: map[string]domain.FCMToken{},
	}
}

func NewMockUserRepositoryWithError(field string, err error) *MockUserRepository {
	m := NewMockUserRepository()
	switch field {
	case "getUser":
		m.getUserErr = err
	case "upsert":
		m.upsertErr = err
	case "upsertFCMToken":
		m.upsertFCMTokenErr = err
	}
	return m
}

func NewMockUserRepositoryEmpty() *MockUserRepository {
	return &MockUserRepository{
		users:      map[string]domain.User{},
		getUserErr: domain.ErrUserNotFound,
		fcmTokens:  map[string]domain.FCMToken{},
	}
}

func (m *MockUserRepository) GetUser(id string) (*domain.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return &user, nil
}

func (m *MockUserRepository) UpsertUser(id string, req domain.UserRequest) (*domain.User, error) {
	if m.upsertErr != nil {
		return nil, m.upsertErr
	}
	user := domain.User{
		ID:     id,
		Email:  req.Email,
		Grade:  req.Grade,
		Course: req.Course,
		Class:  req.Class,
	}
	m.users[id] = user
	return &user, nil
}

func (m *MockUserRepository) UpsertFCMToken(userID string, req domain.FCMTokenRequest) (*domain.FCMToken, error) {
	if m.upsertFCMTokenErr != nil {
		return nil, m.upsertFCMTokenErr
	}

	now := time.Now().UTC()
	token := domain.FCMToken{
		Token:     req.Token,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if existing, ok := m.fcmTokens[userID]; ok {
		token.CreatedAt = existing.CreatedAt
	}
	m.fcmTokens[userID] = token

	return &token, nil
}
