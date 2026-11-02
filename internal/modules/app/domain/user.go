package domain

import (
	"errors"
	"time"
)

// ErrUserNotFound ユーザーが見つからない
var ErrUserNotFound = errors.New("user not found")

// User ユーザー
type User struct {
	ID     string
	Email  string
	Grade  *Grade
	Course *Course
	Class  *Class
}

// UserRequest ユーザー作成・更新リクエスト
type UserRequest struct {
	Email  string
	Grade  *Grade
	Course *Course
	Class  *Class
}

// FCMToken FCMトークン
type FCMToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FCMTokenRequest FCMトークン作成・更新リクエスト
type FCMTokenRequest struct {
	Token string
}
