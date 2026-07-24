package domain

import "time"

type FCMTokenListFilter struct {
	UserIDs       []string
	Tokens        []string
	UpdatedAtFrom *time.Time
	UpdatedAtTo   *time.Time
}
