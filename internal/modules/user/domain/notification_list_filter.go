package domain

import "time"

type NotificationListFilter struct {
	NotifyAtFrom *time.Time
	NotifyAtTo   *time.Time
	IsNotified   *bool
}
