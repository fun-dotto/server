package parser

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

	schedule, err := Parse(directory)
	if err != nil {
		t.Fatalf("parse GTFS models: %v", err)
	}

	if len(schedule.Routes) != 1 || schedule.Routes[0].RouteShortName != "10" {
		t.Fatalf("unexpected routes: %#v", schedule.Routes)
	}
	if len(schedule.Stops) != 1 || !schedule.Stops[0].StopLat.Equal(decimal.RequireFromString("35.1234567")) || schedule.Stops[0].ZoneID == nil || *schedule.Stops[0].ZoneID != "z1" {
		t.Fatalf("unexpected stops: %#v", schedule.Stops)
	}
	if len(schedule.Trips) != 1 || schedule.Trips[0].DirectionID == nil || *schedule.Trips[0].DirectionID != 1 {
		t.Fatalf("unexpected trips: %#v", schedule.Trips)
	}
	if len(schedule.StopTimes) != 1 || schedule.StopTimes[0].StopSequence != 3 || schedule.StopTimes[0].StopHeadsign == nil {
		t.Fatalf("unexpected stop times: %#v", schedule.StopTimes)
	}
	if len(schedule.Calendars) != 1 || !schedule.Calendars[0].StartDate.Equal(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected calendars: %#v", schedule.Calendars)
	}
	if len(schedule.CalendarDates) != 1 || schedule.CalendarDates[0].ExceptionType != 2 {
		t.Fatalf("unexpected calendar dates: %#v", schedule.CalendarDates)
	}
	if len(schedule.FareAttributes) != 1 || !schedule.FareAttributes[0].Price.Equal(decimal.RequireFromString("230.50")) {
		t.Fatalf("unexpected fare attributes: %#v", schedule.FareAttributes)
	}
	if len(schedule.FareRules) != 1 || schedule.FareRules[0].ID == uuid.Nil() || schedule.FareRules[0].OriginID == nil || *schedule.FareRules[0].OriginID != "z2" {
		t.Fatalf("unexpected fare rules: %#v", schedule.FareRules)
	}
	if len(schedule.Zones) != 2 || schedule.Zones[0].ZoneID != "z1" || schedule.Zones[1].ZoneID != "z2" {
		t.Fatalf("unexpected zones: %#v", schedule.Zones)
	}
}

func TestLoadGTFSModelsRejectsInvalidNumber(t *testing.T) {
	directory := t.TempDir()
	contents := "stop_id,stop_name,stop_lat,stop_lon\ns1,Central,not-a-number,139.1\n"
	if err := os.WriteFile(filepath.Join(directory, "stops.txt"), []byte(contents), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := Parse(directory)
	if err == nil {
		t.Fatal("expected invalid stop_lat to fail")
	}
}
