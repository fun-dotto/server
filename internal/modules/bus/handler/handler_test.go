package handler

import (
	"context"
	"testing"
	"time"

	api "github.com/fun-dotto/server/gen/bus"
	"github.com/fun-dotto/server/internal/modules/bus/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type fakeBusService struct{}

func (fakeBusService) ListTripDetails(ctx context.Context, date time.Time) ([]domain.TripDetail, error) {
	return []domain.TripDetail{{
		Trip: domain.Trip{
			TripID:      "trip-1",
			RouteID:     "route-1",
			ServiceID:   "service-1",
			DirectionID: 1,
		},
		Route: domain.Route{RouteID: "route-1", RouteShortName: "1"},
		StopTimes: []domain.StopTimeWithStop{{
			StopTime: domain.StopTime{
				TripID:        "trip-1",
				StopID:        "stop-a",
				DepartureTime: "08:00:00",
				ArrivalTime:   "07:55:00",
				StopSequence:  1,
			},
			Stop: domain.Stop{StopID: "stop-a", StopName: "A"},
		}},
	}}, nil
}

func (fakeBusService) ListTimetableStops(ctx context.Context, tripID string) ([]domain.StopTimeWithStop, error) {
	return []domain.StopTimeWithStop{{
		StopTime: domain.StopTime{
			TripID:        tripID,
			StopID:        "stop-a",
			DepartureTime: "08:00:00",
			ArrivalTime:   "07:55:00",
			StopSequence:  1,
		},
		Stop: domain.Stop{StopID: "stop-a", StopName: "A"},
	}}, nil
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(fakeBusService{})
	if h == nil {
		t.Fatal("handler should be initialized")
	}
}

func TestBusTripsV1List(t *testing.T) {
	h := NewHandler(fakeBusService{})
	resp, err := h.BusTripsV1List(context.Background(), api.BusTripsV1ListRequestObject{
		Params: api.BusTripsV1ListParams{Date: openapi_types.Date{Time: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonResp, ok := resp.(api.BusTripsV1List200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if len(jsonResp.BusTrips) != 1 {
		t.Fatalf("expected 1 trip, got %d", len(jsonResp.BusTrips))
	}
	if jsonResp.BusTrips[0].Route.Name != "1" {
		t.Fatalf("expected route name 1, got %s", jsonResp.BusTrips[0].Route.Name)
	}
}

func TestBusTimetableStopsV1List(t *testing.T) {
	h := NewHandler(fakeBusService{})
	resp, err := h.BusTimetableStopsV1List(context.Background(), api.BusTimetableStopsV1ListRequestObject{TripId: "trip-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonResp, ok := resp.(api.BusTimetableStopsV1List200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if len(jsonResp.BusTimetableStops) != 1 {
		t.Fatalf("expected 1 stop, got %d", len(jsonResp.BusTimetableStops))
	}
	if jsonResp.BusTimetableStops[0].Stop.Name != "A" {
		t.Fatalf("expected stop name A, got %s", jsonResp.BusTimetableStops[0].Stop.Name)
	}
}
