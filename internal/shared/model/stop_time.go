package model

type StopTime struct {
	TripID 	      string `gorm:"primaryKey;index"`
	Trip          *Trip `gorm:"belongsTo;foreignKey:TripID;references:TripID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ArrivalTime   *string
	DepartureTime *string
	StopID        string `gorm:"not null;index"`
	Stop          *Stop  `gorm:"belongsTo;foreignKey:StopID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StopSequence  int    `gorm:"primaryKey"`
}
