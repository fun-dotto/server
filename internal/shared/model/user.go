package model

type User struct {
	ID    string `gorm:"primaryKey"`
	Email string `gorm:"not null"`

	Grade  *string
	Course *string
	Class  *string
}
