package model

type Calendar struct {
	Common

	ServiceID string `gorm:"not null;uniqueIndex"`
	Monday    int    `gorm:"not null"`
	Tuesday   int    `gorm:"not null"`
	Wednesday int    `gorm:"not null"`
	Thursday  int    `gorm:"not null"`
	Friday    int    `gorm:"not null"`
	Saturday  int    `gorm:"not null"`
	Sunday    int    `gorm:"not null"`
	StartDate string `gorm:"not null"`
	EndDate   string `gorm:"not null"`
}
