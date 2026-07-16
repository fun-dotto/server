package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func parseStops(files map[string]*zip.File) ([]Stop, error) {
	rows, err := readCSVFile(files, "stops.txt")
	if err != nil {
		return nil, err
	}
	out := make([]Stop, 0, len(rows))
	for _, row := range rows {
		out = append(out, Stop{
			StopID:             row["stop_id"],
			StopCode:           row["stop_code"],
			StopName:           row["stop_name"],
			StopDesc:           row["stop_desc"],
			StopLat:            row["stop_lat"],
			StopLon:            row["stop_lon"],
			ZoneID:             row["zone_id"],
			StopURL:            row["stop_url"],
			LocationType:       row["location_type"],
			ParentStation:      row["parent_station"],
			StopTimezone:       row["stop_timezone"],
			WheelchairBoarding: row["wheelchair_boarding"],
		})
	}
	return out, nil
}

func parseRoutes(files map[string]*zip.File) ([]Route, error) {
	rows, err := readCSVFile(files, "routes.txt")
	if err != nil {
		return nil, err
	}
	out := make([]Route, 0, len(rows))
	for _, row := range rows {
		out = append(out, Route{
			RouteID:         row["route_id"],
			AgencyID:        row["agency_id"],
			RouteShortName:  row["route_short_name"],
			RouteLongName:   row["route_long_name"],
			RouteDesc:       row["route_desc"],
			RouteType:       row["route_type"],
			RouteURL:        row["route_url"],
			RouteColor:      row["route_color"],
			RouteTextColor:  row["route_text_color"],
			JPParentRouteID: row["jp_parent_route_id"],
		})
	}
	return out, nil
}

func parseCalendars(files map[string]*zip.File) ([]Calendar, error) {
	rows, err := readCSVFile(files, "calendar.txt")
	if err != nil {
		return nil, err
	}
	out := make([]Calendar, 0, len(rows))
	for _, row := range rows {
		out = append(out, Calendar{
			ServiceID: row["service_id"],
			Monday:    row["monday"],
			Tuesday:   row["tuesday"],
			Wednesday: row["wednesday"],
			Thursday:  row["thursday"],
			Friday:    row["friday"],
			Saturday:  row["saturday"],
			Sunday:    row["sunday"],
			StartDate: row["start_date"],
			EndDate:   row["end_date"],
		})
	}
	return out, nil
}

func parseCalendarDates(files map[string]*zip.File) ([]CalendarDate, error) {
	rows, err := readCSVFile(files, "calendar_dates.txt")
	if err != nil {
		return nil, err
	}
	out := make([]CalendarDate, 0, len(rows))
	for _, row := range rows {
		out = append(out, CalendarDate{
			ServiceID:     row["service_id"],
			Date:          row["date"],
			ExceptionType: row["exception_type"],
		})
	}
	return out, nil
}

func parseTrips(files map[string]*zip.File) ([]Trip, error) {
	rows, err := readCSVFile(files, "trips.txt")
	if err != nil {
		return nil, err
	}
	out := make([]Trip, 0, len(rows))
	for _, row := range rows {
		out = append(out, Trip{
			RouteID:              row["route_id"],
			ServiceID:            row["service_id"],
			TripID:               row["trip_id"],
			TripHeadsign:         row["trip_headsign"],
			TripShortName:        row["trip_short_name"],
			DirectionID:          row["direction_id"],
			BlockID:              row["block_id"],
			ShapeID:              row["shape_id"],
			WheelchairAccessible: row["wheelchair_accessible"],
			BikesAllowed:         row["bikes_allowed"],
			JPTripDesc:           row["jp_trip_desc"],
			JPTripDescSymbol:     row["jp_trip_desc_symbol"],
			JPOfficeID:           row["jp_office_id"],
		})
	}
	return out, nil
}

func parseStopTimes(files map[string]*zip.File) ([]StopTime, error) {
	rows, err := readCSVFile(files, "stop_times.txt")
	if err != nil {
		return nil, err
	}
	out := make([]StopTime, 0, len(rows))
	for _, row := range rows {
		out = append(out, StopTime{
			TripID:            row["trip_id"],
			ArrivalTime:       row["arrival_time"],
			DepartureTime:     row["departure_time"],
			StopID:            row["stop_id"],
			StopSequence:      row["stop_sequence"],
			StopHeadsign:      row["stop_headsign"],
			PickupType:        row["pickup_type"],
			DropOffType:       row["drop_off_type"],
			ShapeDistTraveled: row["shape_dist_traveled"],
			Timepoint:         row["timepoint"],
		})
	}
	return out, nil
}

func parseFareTables(files map[string]*zip.File) ([]FareRule, []FareAttribute, error) {
	ruleRows, err := readCSVFile(files, "fare_rules.txt")
	if err != nil {
		return nil, nil, err
	}
	attrRows, err := readCSVFile(files, "fare_attributes.txt")
	if err != nil {
		return nil, nil, err
	}

	rules := make([]FareRule, 0, len(ruleRows))
	for _, row := range ruleRows {
		rules = append(rules, FareRule{
			FareID:        row["fare_id"],
			RouteID:       row["route_id"],
			OriginID:      row["origin_id"],
			DestinationID: row["destination_id"],
			ContainsID:    row["contains_id"],
		})
	}

	attrs := make([]FareAttribute, 0, len(attrRows))
	for _, row := range attrRows {
		attrs = append(attrs, FareAttribute{
			FareID:           row["fare_id"],
			Price:            row["price"],
			CurrencyType:     row["currency_type"],
			PaymentMethod:    row["payment_method"],
			Transfers:        row["transfers"],
			TransferDuration: row["transfer_duration"],
		})
	}

	return rules, attrs, nil
}

func readCSVFile(files map[string]*zip.File, name string) ([]map[string]string, error) {
	f, ok := files[name]
	if !ok {
		return nil, fmt.Errorf("%s not found in zip", name)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	defer rc.Close()

	rows, err := readCSVMaps(rc)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	return rows, nil
}

func readCSVMaps(r io.Reader) ([]map[string]string, error) {
	cr := csv.NewReader(r)
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	var rows []map[string]string
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		row := make(map[string]string, len(header))
		for i, key := range header {
			if i < len(record) {
				row[key] = record[i]
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
