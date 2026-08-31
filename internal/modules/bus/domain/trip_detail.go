package domain

// TripDetail は運行情報に路線・停車列を組み合わせた集約。
type TripDetail struct {
	Trip      Trip
	Route     Route
	StopTimes []StopTimeWithStop
}
