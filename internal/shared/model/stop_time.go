package model

type StopTime struct {
	TripID string `gorm:"primaryKey"`
	// belongsTo は非公式だが gorm.io/gorm の schema/relationship.go で実際に処理されるタグなので削除しない
	Trip          *Trip `gorm:"belongsTo;foreignKey:TripID;references:TripID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ArrivalTime   *string
	DepartureTime *string
	StopID        string `gorm:"not null;index"`
	Stop          *Stop  `gorm:"belongsTo;foreignKey:StopID;references:StopID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StopSequence  int    `gorm:"primaryKey"`
}
