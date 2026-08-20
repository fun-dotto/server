package model

type StopTime struct {
	Common

	TripID        string `gorm:"not null;index"`
	Trip          *Trip  `gorm:"belongsTo;foreignKey:TripID;references:TripID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ArrivalTime   string `gorm:"not null"`
	DepartureTime string `gorm:"not null"`
	StopID        string `gorm:"not null;index"`
	Stop          *Stop  `gorm:"belongsTo;foreignKey:StopID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StopSequence  int    `gorm:"not null"`
}
