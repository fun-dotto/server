package domain

// StopTimeWithStop は停車時刻とその停留所情報を組み合わせた集約。
type StopTimeWithStop struct {
	StopTime StopTime
	Stop     Stop
}
