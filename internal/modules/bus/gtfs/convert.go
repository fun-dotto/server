package gtfs

import (
	"fmt"
	"strconv"

	"github.com/fun-dotto/server/internal/modules/bus/domain"
)

func ToDomainStop(row Stop) (domain.Stop, error) {
	if row.StopID == "" {
		return domain.Stop{}, fmt.Errorf("stop_id is required")
	}
	if row.StopName == "" {
		return domain.Stop{}, fmt.Errorf("stop_name is required for stop_id=%s", row.StopID)
	}
	return domain.Stop{
		StopID:   row.StopID,
		StopName: row.StopName,
	}, nil
}

func ToDomainStops(rows []Stop) ([]domain.Stop, error) {
	out := make([]domain.Stop, 0, len(rows))
	for i, row := range rows {
		stop, err := ToDomainStop(row)
		if err != nil {
			return nil, fmt.Errorf("stops.txt row %d: %w", i+2, err)
		}
		out = append(out, stop)
	}
	return out, nil
}

func ToDomainRoute(row Route) (domain.Route, error) {
	if row.RouteID == "" {
		return domain.Route{}, fmt.Errorf("route_id is required")
	}
	return domain.Route{
		RouteID:        row.RouteID,
		RouteShortName: row.RouteShortName,
	}, nil
}

func ToDomainRoutes(rows []Route) ([]domain.Route, error) {
	out := make([]domain.Route, 0, len(rows))
	for i, row := range rows {
		route, err := ToDomainRoute(row)
		if err != nil {
			return nil, fmt.Errorf("routes.txt row %d: %w", i+2, err)
		}
		out = append(out, route)
	}
	return out, nil
}

func ToDomainCalendar(row Calendar) (domain.Calendar, error) {
	if row.ServiceID == "" {
		return domain.Calendar{}, fmt.Errorf("service_id is required")
	}
	monday, err := parseInt(row.Monday, "monday")
	if err != nil {
		return domain.Calendar{}, err
	}
	tuesday, err := parseInt(row.Tuesday, "tuesday")
	if err != nil {
		return domain.Calendar{}, err
	}
	wednesday, err := parseInt(row.Wednesday, "wednesday")
	if err != nil {
		return domain.Calendar{}, err
	}
	thursday, err := parseInt(row.Thursday, "thursday")
	if err != nil {
		return domain.Calendar{}, err
	}
	friday, err := parseInt(row.Friday, "friday")
	if err != nil {
		return domain.Calendar{}, err
	}
	saturday, err := parseInt(row.Saturday, "saturday")
	if err != nil {
		return domain.Calendar{}, err
	}
	sunday, err := parseInt(row.Sunday, "sunday")
	if err != nil {
		return domain.Calendar{}, err
	}
	if row.StartDate == "" || row.EndDate == "" {
		return domain.Calendar{}, fmt.Errorf("start_date and end_date are required for service_id=%s", row.ServiceID)
	}
	return domain.Calendar{
		ServiceID: row.ServiceID,
		Monday:    monday,
		Tuesday:   tuesday,
		Wednesday: wednesday,
		Thursday:  thursday,
		Friday:    friday,
		Saturday:  saturday,
		Sunday:    sunday,
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
	}, nil
}

func ToDomainCalendars(rows []Calendar) ([]domain.Calendar, error) {
	out := make([]domain.Calendar, 0, len(rows))
	for i, row := range rows {
		calendar, err := ToDomainCalendar(row)
		if err != nil {
			return nil, fmt.Errorf("calendar.txt row %d: %w", i+2, err)
		}
		out = append(out, calendar)
	}
	return out, nil
}

func ToDomainCalendarDate(row CalendarDate) (domain.CalendarDate, error) {
	if row.ServiceID == "" {
		return domain.CalendarDate{}, fmt.Errorf("service_id is required")
	}
	if row.Date == "" {
		return domain.CalendarDate{}, fmt.Errorf("date is required for service_id=%s", row.ServiceID)
	}
	exceptionType, err := parseInt(row.ExceptionType, "exception_type")
	if err != nil {
		return domain.CalendarDate{}, err
	}
	return domain.CalendarDate{
		ServiceID:     row.ServiceID,
		Date:          row.Date,
		ExceptionType: exceptionType,
	}, nil
}

func ToDomainCalendarDates(rows []CalendarDate) ([]domain.CalendarDate, error) {
	out := make([]domain.CalendarDate, 0, len(rows))
	for i, row := range rows {
		calendarDate, err := ToDomainCalendarDate(row)
		if err != nil {
			return nil, fmt.Errorf("calendar_dates.txt row %d: %w", i+2, err)
		}
		out = append(out, calendarDate)
	}
	return out, nil
}

func ToDomainTrip(row Trip) (domain.Trip, error) {
	if row.TripID == "" {
		return domain.Trip{}, fmt.Errorf("trip_id is required")
	}
	if row.RouteID == "" {
		return domain.Trip{}, fmt.Errorf("route_id is required for trip_id=%s", row.TripID)
	}
	if row.ServiceID == "" {
		return domain.Trip{}, fmt.Errorf("service_id is required for trip_id=%s", row.TripID)
	}
	directionID, err := parseInt(row.DirectionID, "direction_id")
	if err != nil {
		return domain.Trip{}, fmt.Errorf("trip_id=%s: %w", row.TripID, err)
	}
	return domain.Trip{
		TripID:      row.TripID,
		RouteID:     row.RouteID,
		ServiceID:   row.ServiceID,
		DirectionID: directionID,
	}, nil
}

func ToDomainTrips(rows []Trip) ([]domain.Trip, error) {
	out := make([]domain.Trip, 0, len(rows))
	for i, row := range rows {
		trip, err := ToDomainTrip(row)
		if err != nil {
			return nil, fmt.Errorf("trips.txt row %d: %w", i+2, err)
		}
		out = append(out, trip)
	}
	return out, nil
}

func ToDomainStopTime(row StopTime) (domain.StopTime, error) {
	if row.TripID == "" {
		return domain.StopTime{}, fmt.Errorf("trip_id is required")
	}
	if row.StopID == "" {
		return domain.StopTime{}, fmt.Errorf("stop_id is required for trip_id=%s", row.TripID)
	}
	if row.ArrivalTime == "" || row.DepartureTime == "" {
		return domain.StopTime{}, fmt.Errorf("arrival_time and departure_time are required for trip_id=%s stop_id=%s", row.TripID, row.StopID)
	}
	stopSequence, err := parseInt(row.StopSequence, "stop_sequence")
	if err != nil {
		return domain.StopTime{}, fmt.Errorf("trip_id=%s stop_id=%s: %w", row.TripID, row.StopID, err)
	}
	return domain.StopTime{
		TripID:        row.TripID,
		ArrivalTime:   row.ArrivalTime,
		DepartureTime: row.DepartureTime,
		StopID:        row.StopID,
		StopSequence:  stopSequence,
	}, nil
}

func ToDomainStopTimes(rows []StopTime) ([]domain.StopTime, error) {
	out := make([]domain.StopTime, 0, len(rows))
	for i, row := range rows {
		stopTime, err := ToDomainStopTime(row)
		if err != nil {
			return nil, fmt.Errorf("stop_times.txt row %d: %w", i+2, err)
		}
		out = append(out, stopTime)
	}
	return out, nil
}

// ToDomainFareRules は fare_rules と fare_attributes を結合して domain にする。
// origin_id / destination_id が欠ける行はスキップする。
func ToDomainFareRules(rules []FareRule, attrs []FareAttribute) ([]domain.FareRule, error) {
	priceByFareID := make(map[string]float64, len(attrs))
	for _, attr := range attrs {
		if attr.FareID == "" {
			return nil, fmt.Errorf("fare_attributes: fare_id is required")
		}
		price, err := strconv.ParseFloat(attr.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("fare_attributes fare_id=%s: invalid price %q: %w", attr.FareID, attr.Price, err)
		}
		priceByFareID[attr.FareID] = price
	}

	out := make([]domain.FareRule, 0, len(rules))
	for _, rule := range rules {
		if rule.FareID == "" {
			return nil, fmt.Errorf("fare_rules: fare_id is required")
		}
		if rule.RouteID == "" || rule.OriginID == "" || rule.DestinationID == "" {
			continue
		}
		price, ok := priceByFareID[rule.FareID]
		if !ok {
			return nil, fmt.Errorf("fare_rules fare_id=%s: matching fare_attributes not found", rule.FareID)
		}
		out = append(out, domain.FareRule{
			RouteID:       rule.RouteID,
			OriginID:      rule.OriginID,
			DestinationID: rule.DestinationID,
			Price:         price,
		})
	}
	return out, nil
}

func parseInt(raw, field string) (int, error) {
	if raw == "" {
		return 0, fmt.Errorf("%s is required", field)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", field, raw, err)
	}
	return v, nil
}
