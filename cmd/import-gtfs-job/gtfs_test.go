package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
)

func TestLoadGTFSModels(t *testing.T) {
	directory := t.TempDir()
	fixtures := map[string]string{
		"routes.txt":          "route_id,route_short_name\nr1,10\n",
		"stops.txt":           "stop_id,stop_name,stop_lat,stop_lon,zone_id\ns1,Central,35.1234567,139.1234567,z1\n",
		"trips.txt":           "route_id,service_id,trip_id,direction_id\nr1,weekday,t1,1\n",
		"stop_times.txt":      "trip_id,stop_sequence,stop_id,arrival_time,departure_time,stop_headsign\nt1,3,s1,08:00:00,08:00:00,Central\n",
		"calendar.txt":        "service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date\nweekday,1,1,1,1,1,0,0,20260101,20261231\n",
		"calendar_dates.txt":  "service_id,date,exception_type\nweekday,20260102,2\n",
		"fare_attributes.txt": "fare_id,price\nf1,230.50\n",
		"fare_rules.txt":      "fare_id,route_id,origin_id,destination_id\nf1,r1,z2,z1\n",
	}
	for filename, contents := range fixtures {
		if err := os.WriteFile(filepath.Join(directory, filename), []byte(contents), 0600); err != nil {
			t.Fatalf("write fixture %s: %v", filename, err)
		}
	}

	models, err := loadGTFSModels(directory)
	if err != nil {
		t.Fatalf("load GTFS models: %v", err)
	}

	if len(models.routes) != 1 || models.routes[0].RouteShortName != "10" {
		t.Fatalf("unexpected routes: %#v", models.routes)
	}
	if len(models.stops) != 1 || !models.stops[0].StopLat.Equal(decimal.RequireFromString("35.1234567")) || models.stops[0].ZoneID == nil || *models.stops[0].ZoneID != "z1" {
		t.Fatalf("unexpected stops: %#v", models.stops)
	}
	if len(models.trips) != 1 || models.trips[0].DirectionID == nil || *models.trips[0].DirectionID != 1 {
		t.Fatalf("unexpected trips: %#v", models.trips)
	}
	if len(models.stopTimes) != 1 || models.stopTimes[0].StopSequence != 3 || models.stopTimes[0].StopHeadsign == nil {
		t.Fatalf("unexpected stop times: %#v", models.stopTimes)
	}
	if len(models.calendars) != 1 || !models.calendars[0].StartDate.Equal(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected calendars: %#v", models.calendars)
	}
	if len(models.calendarDates) != 1 || models.calendarDates[0].ExceptionType != 2 {
		t.Fatalf("unexpected calendar dates: %#v", models.calendarDates)
	}
	if len(models.fareAttributes) != 1 || !models.fareAttributes[0].Price.Equal(decimal.RequireFromString("230.50")) {
		t.Fatalf("unexpected fare attributes: %#v", models.fareAttributes)
	}
	if len(models.fareRules) != 1 || models.fareRules[0].ID == uuid.Nil() || models.fareRules[0].OriginID == nil || *models.fareRules[0].OriginID != "z2" {
		t.Fatalf("unexpected fare rules: %#v", models.fareRules)
	}
	if len(models.zones) != 2 || models.zones[0].ZoneID != "z1" || models.zones[1].ZoneID != "z2" {
		t.Fatalf("unexpected zones: %#v", models.zones)
	}
}

func TestLoadGTFSModelsRejectsInvalidNumber(t *testing.T) {
	directory := t.TempDir()
	contents := "stop_id,stop_name,stop_lat,stop_lon\ns1,Central,not-a-number,139.1\n"
	if err := os.WriteFile(filepath.Join(directory, "stops.txt"), []byte(contents), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := loadGTFSModels(directory)
	if err == nil {
		t.Fatal("expected invalid stop_lat to fail")
	}
}
