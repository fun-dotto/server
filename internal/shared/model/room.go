package model

type Room struct {
	Common

	Name  string `gorm:"not null"`
	Floor string `gorm:"not null"`
}
