package model

type StopTime struct {
	TripID        string  `gorm:"primaryKey"`
	Trip          *Trip   `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
	StopSequence  int     `gorm:"primaryKey"`
	StopID        string  `gorm:"not null;index"`
	Stop          *Stop   `gorm:"foreignKey:StopID;constraint:OnDelete:CASCADE"`
	ArrivalTime   string  `gorm:"not null"`
	DepartureTime string  `gorm:"not null"`
	StopHeadsign  *string `gorm:"type:text"`
}
