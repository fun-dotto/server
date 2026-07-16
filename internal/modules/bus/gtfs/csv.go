package gtfs

// GTFS-JP の実ファイルに忠実な CSV モデル。
// 1 つでもプロパティを使うファイルのみを定義し、各ファイルの全カラムを持つ。
// domain への整形（不要カラムの間引き・型変換）は convert.go 側で行う。

// Stop は stops.txt。
type Stop struct {
	StopID             string // stop_id
	StopCode           string // stop_code
	StopName           string // stop_name
	StopDesc           string // stop_desc
	StopLat            string // stop_lat
	StopLon            string // stop_lon
	ZoneID             string // zone_id
	StopURL            string // stop_url
	LocationType       string // location_type
	ParentStation      string // parent_station
	StopTimezone       string // stop_timezone
	WheelchairBoarding string // wheelchair_boarding
}

// Route は routes.txt。
type Route struct {
	RouteID         string // route_id
	AgencyID        string // agency_id
	RouteShortName  string // route_short_name
	RouteLongName   string // route_long_name
	RouteDesc       string // route_desc
	RouteType       string // route_type
	RouteURL        string // route_url
	RouteColor      string // route_color
	RouteTextColor  string // route_text_color
	JPParentRouteID string // jp_parent_route_id
}

// Calendar は calendar.txt。
type Calendar struct {
	ServiceID string // service_id
	Monday    string // monday
	Tuesday   string // tuesday
	Wednesday string // wednesday
	Thursday  string // thursday
	Friday    string // friday
	Saturday  string // saturday
	Sunday    string // sunday
	StartDate string // start_date
	EndDate   string // end_date
}

// CalendarDate は calendar_dates.txt。
type CalendarDate struct {
	ServiceID     string // service_id
	Date          string // date
	ExceptionType string // exception_type
}

// Trip は trips.txt。
type Trip struct {
	RouteID              string // route_id
	ServiceID            string // service_id
	TripID               string // trip_id
	TripHeadsign         string // trip_headsign
	TripShortName        string // trip_short_name
	DirectionID          string // direction_id
	BlockID              string // block_id
	ShapeID              string // shape_id
	WheelchairAccessible string // wheelchair_accessible
	BikesAllowed         string // bikes_allowed
	JPTripDesc           string // jp_trip_desc
	JPTripDescSymbol     string // jp_trip_desc_symbol
	JPOfficeID           string // jp_office_id
}

// StopTime は stop_times.txt。
type StopTime struct {
	TripID            string // trip_id
	ArrivalTime       string // arrival_time
	DepartureTime     string // departure_time
	StopID            string // stop_id
	StopSequence      string // stop_sequence
	StopHeadsign      string // stop_headsign
	PickupType        string // pickup_type
	DropOffType       string // drop_off_type
	ShapeDistTraveled string // shape_dist_traveled
	Timepoint         string // timepoint
}

// FareRule は fare_rules.txt。
type FareRule struct {
	FareID        string // fare_id
	RouteID       string // route_id
	OriginID      string // origin_id
	DestinationID string // destination_id
	ContainsID    string // contains_id
}

// FareAttribute は fare_attributes.txt。
type FareAttribute struct {
	FareID           string // fare_id
	Price            string // price
	CurrencyType     string // currency_type
	PaymentMethod    string // payment_method
	Transfers        string // transfers
	TransferDuration string // transfer_duration
}
