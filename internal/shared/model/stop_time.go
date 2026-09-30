package model

type StopTime struct {
	TripID        string `gorm:"primaryKey"`
	ArrivalTime   string `gorm:"not null"`
	DepartureTime string `gorm:"not null"`
	StopID        string `gorm:"not null;index"`
	StopSequence  int    `gorm:"primaryKey"`
	Trip          Trip   `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
	Stop          Stop   `gorm:"foreignKey:StopID;constraint:OnDelete:CASCADE"`
}
