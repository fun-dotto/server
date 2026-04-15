package database

type NotificationTargetUser struct {
	NotificationID string       `gorm:"type:text;primaryKey"`
	UserID         string       `gorm:"type:text;primaryKey"`
	Notification   Notification `gorm:"constraint:OnDelete:CASCADE"`
	User           User         `gorm:"constraint:OnDelete:CASCADE"`
}
