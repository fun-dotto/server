package model

type StopTime struct {
	TripID        string  `gorm:"primaryKey"`
	StopSequence  int     `gorm:"primaryKey"`
	StopID        string  `gorm:"not null;index"`
	ArrivalTime   string  `gorm:"not null"`
	DepartureTime string  `gorm:"not null"`
	StopHeadsign  *string `gorm:"type:text"`
	Trip          Trip    `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
	Stop          Stop    `gorm:"foreignKey:StopID;constraint:OnDelete:CASCADE"`
}
