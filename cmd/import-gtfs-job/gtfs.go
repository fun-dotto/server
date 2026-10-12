package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"uuid"

	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/shopspring/decimal"
)

type gtfsModels struct {
	calendars      []model.Calendar
	calendarDates  []model.CalendarDate
	zones          []model.Zone
	fareAttributes []model.FareAttribute
	fareRules      []model.FareRule
	routes         []model.Route
	stops          []model.Stop
	trips          []model.Trip
	stopTimes      []model.StopTime
}

func loadGTFSModels(directory string) (gtfsModels, error) {
	var models gtfsModels
	zoneIDs := make(map[string]struct{})

	if err := readGTFSCSV(directory, "routes.txt", []string{"route_id"}, func(row map[string]string, line int) error {
		routeID, err := requiredGTFSValue(row, "route_id", "routes.txt", line)
		if err != nil {
			return err
		}
		models.routes = append(models.routes, model.Route{
			RouteID:        routeID,
			RouteShortName: row["route_short_name"],
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "stops.txt", []string{"stop_id", "stop_name", "stop_lat", "stop_lon"}, func(row map[string]string, line int) error {
		stopID, err := requiredGTFSValue(row, "stop_id", "stops.txt", line)
		if err != nil {
			return err
		}
		stopName, err := requiredGTFSValue(row, "stop_name", "stops.txt", line)
		if err != nil {
			return err
		}
		latitude, err := decimal.NewFromString(row["stop_lat"])
		if err != nil {
			return fmt.Errorf("stops.txt line %d: invalid stop_lat %q: %w", line, row["stop_lat"], err)
		}
		longitude, err := decimal.NewFromString(row["stop_lon"])
		if err != nil {
			return fmt.Errorf("stops.txt line %d: invalid stop_lon %q: %w", line, row["stop_lon"], err)
		}
		zoneID := optionalGTFSValue(row["zone_id"])
		if zoneID != nil {
			zoneIDs[*zoneID] = struct{}{}
		}
		models.stops = append(models.stops, model.Stop{
			StopID:   stopID,
			StopName: stopName,
			StopLat:  latitude,
			StopLon:  longitude,
			ZoneID:   zoneID,
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "trips.txt", []string{"trip_id", "route_id", "service_id"}, func(row map[string]string, line int) error {
		tripID, err := requiredGTFSValue(row, "trip_id", "trips.txt", line)
		if err != nil {
			return err
		}
		routeID, err := requiredGTFSValue(row, "route_id", "trips.txt", line)
		if err != nil {
			return err
		}
		serviceID, err := requiredGTFSValue(row, "service_id", "trips.txt", line)
		if err != nil {
			return err
		}
		var directionID *int
		if value := strings.TrimSpace(row["direction_id"]); value != "" {
			direction, err := strconv.Atoi(value)
			if err != nil || (direction != 0 && direction != 1) {
				return fmt.Errorf("trips.txt line %d: invalid direction_id %q", line, value)
			}
			directionID = &direction
		}
		models.trips = append(models.trips, model.Trip{
			TripID:      tripID,
			RouteID:     routeID,
			ServiceID:   serviceID,
			DirectionID: directionID,
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "stop_times.txt", []string{"trip_id", "stop_sequence", "stop_id", "arrival_time", "departure_time"}, func(row map[string]string, line int) error {
		tripID, err := requiredGTFSValue(row, "trip_id", "stop_times.txt", line)
		if err != nil {
			return err
		}
		stopID, err := requiredGTFSValue(row, "stop_id", "stop_times.txt", line)
		if err != nil {
			return err
		}
		sequence, err := strconv.Atoi(row["stop_sequence"])
		if err != nil {
			return fmt.Errorf("stop_times.txt line %d: invalid stop_sequence %q: %w", line, row["stop_sequence"], err)
		}
		models.stopTimes = append(models.stopTimes, model.StopTime{
			TripID:        tripID,
			StopSequence:  sequence,
			StopID:        stopID,
			ArrivalTime:   row["arrival_time"],
			DepartureTime: row["departure_time"],
			StopHeadsign:  optionalGTFSValue(row["stop_headsign"]),
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "calendar.txt", []string{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"}, func(row map[string]string, line int) error {
		serviceID, err := requiredGTFSValue(row, "service_id", "calendar.txt", line)
		if err != nil {
			return err
		}
		weekdays := make([]bool, 7)
		for index, field := range []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"} {
			value := row[field]
			if value != "0" && value != "1" {
				return fmt.Errorf("calendar.txt line %d: invalid %s %q; expected 0 or 1", line, field, value)
			}
			weekdays[index] = value == "1"
		}
		startDate, err := parseGTFSDate(row["start_date"], "calendar.txt", "start_date", line)
		if err != nil {
			return err
		}
		endDate, err := parseGTFSDate(row["end_date"], "calendar.txt", "end_date", line)
		if err != nil {
			return err
		}
		models.calendars = append(models.calendars, model.Calendar{
			ServiceID: serviceID,
			Monday:    weekdays[0], Tuesday: weekdays[1], Wednesday: weekdays[2], Thursday: weekdays[3],
			Friday: weekdays[4], Saturday: weekdays[5], Sunday: weekdays[6],
			StartDate: startDate, EndDate: endDate,
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "calendar_dates.txt", []string{"service_id", "date", "exception_type"}, func(row map[string]string, line int) error {
		serviceID, err := requiredGTFSValue(row, "service_id", "calendar_dates.txt", line)
		if err != nil {
			return err
		}
		date, err := parseGTFSDate(row["date"], "calendar_dates.txt", "date", line)
		if err != nil {
			return err
		}
		exceptionType, err := strconv.Atoi(row["exception_type"])
		if err != nil || (exceptionType != 1 && exceptionType != 2) {
			return fmt.Errorf("calendar_dates.txt line %d: invalid exception_type %q", line, row["exception_type"])
		}
		models.calendarDates = append(models.calendarDates, model.CalendarDate{
			ServiceID: serviceID, Date: date, ExceptionType: exceptionType,
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "fare_attributes.txt", []string{"fare_id", "price"}, func(row map[string]string, line int) error {
		fareID, err := requiredGTFSValue(row, "fare_id", "fare_attributes.txt", line)
		if err != nil {
			return err
		}
		price, err := decimal.NewFromString(row["price"])
		if err != nil {
			return fmt.Errorf("fare_attributes.txt line %d: invalid price %q: %w", line, row["price"], err)
		}
		models.fareAttributes = append(models.fareAttributes, model.FareAttribute{
			FareID: fareID, RiderCategory: model.RiderCategoryAdult, Price: price,
		})
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	if err := readGTFSCSV(directory, "fare_rules.txt", []string{"fare_id"}, func(row map[string]string, line int) error {
		fareID, err := requiredGTFSValue(row, "fare_id", "fare_rules.txt", line)
		if err != nil {
			return err
		}
		rule := model.FareRule{
			ID:            uuid.New(),
			FareID:        fareID,
			RouteID:       optionalGTFSValue(row["route_id"]),
			OriginID:      optionalGTFSValue(row["origin_id"]),
			DestinationID: optionalGTFSValue(row["destination_id"]),
		}
		if rule.OriginID != nil {
			zoneIDs[*rule.OriginID] = struct{}{}
		}
		if rule.DestinationID != nil {
			zoneIDs[*rule.DestinationID] = struct{}{}
		}
		models.fareRules = append(models.fareRules, rule)
		return nil
	}); err != nil {
		return gtfsModels{}, err
	}

	orderedZoneIDs := make([]string, 0, len(zoneIDs))
	for zoneID := range zoneIDs {
		orderedZoneIDs = append(orderedZoneIDs, zoneID)
	}
	sort.Strings(orderedZoneIDs)
	for _, zoneID := range orderedZoneIDs {
		models.zones = append(models.zones, model.Zone{ZoneID: zoneID})
	}

	return models, nil
}

func readGTFSCSV(directory, filename string, requiredHeaders []string, consume func(map[string]string, int) error) error {
	path := filepath.Join(directory, filename)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read %s header: %w", filename, err)
	}
	if len(headers) > 0 {
		headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
	}
	headerIndexes := make(map[string]int, len(headers))
	for index, header := range headers {
		headerIndexes[strings.TrimSpace(header)] = index
	}
	for _, header := range requiredHeaders {
		if _, ok := headerIndexes[header]; !ok {
			return fmt.Errorf("%s is missing required column %q", filename, header)
		}
	}

	line := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read %s line %d: %w", filename, line, err)
		}
		row := make(map[string]string, len(headers))
		for name, index := range headerIndexes {
			row[name] = strings.TrimSpace(record[index])
		}
		if err := consume(row, line); err != nil {
			return err
		}
		line++
	}
}

func requiredGTFSValue(row map[string]string, field, filename string, line int) (string, error) {
	value := strings.TrimSpace(row[field])
	if value == "" {
		return "", fmt.Errorf("%s line %d: %s is empty", filename, line, field)
	}
	return value, nil
}

func optionalGTFSValue(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func parseGTFSDate(value, filename, field string, line int) (time.Time, error) {
	date, err := time.Parse("20060102", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s line %d: invalid %s %q: %w", filename, line, field, value, err)
	}
	return date, nil
}
