package domain

type FareRule struct {
	ID            string
	RouteID       string
	OriginID      string
	DestinationID string
	Price         float64
}
