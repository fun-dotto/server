package domain

import "time"

type FCMToken struct {
	Token     string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
