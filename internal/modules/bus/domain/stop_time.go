package domain

type StopTime struct {
	ID            string
	TripID        string
	ArrivalTime   string
	DepartureTime string
	StopID        string
	StopSequence  int
}
