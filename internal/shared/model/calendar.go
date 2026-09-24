package model

import "time"

type Calendar struct {
	Common

	ServiceID string    `gorm:"not null;uniqueIndex"`
	Monday    int       `gorm:"not null;check:monday IN (0,1)"`
	Tuesday   int       `gorm:"not null;check:tuesday IN (0,1)"`
	Wednesday int       `gorm:"not null;check:wednesday IN (0,1)"`
	Thursday  int       `gorm:"not null;check:thursday IN (0,1)"`
	Friday    int       `gorm:"not null;check:friday IN (0,1)"`
	Saturday  int       `gorm:"not null;check:saturday IN (0,1)"`
	Sunday    int       `gorm:"not null;check:sunday IN (0,1)"`
	StartDate time.Time `gorm:"type:date;not null"`
	EndDate   time.Time `gorm:"type:date;not null"`
}
