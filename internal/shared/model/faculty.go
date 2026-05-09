package model

type Faculty struct {
	Common

	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}
