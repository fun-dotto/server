package repository

import (
	"github.com/fun-dotto/server/internal/modules/bus/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

// domain（bus）→ shared/model の境界変換。
// 内部 ID（uuid）は DB 側の default 採番に任せ、業務キー（stop_id 等）のみ設定する。

func stopFromDomain(d domain.Stop) model.Stop {
	return model.Stop{
		StopID:   d.StopID,
		StopName: d.StopName,
	}
}

func routeFromDomain(d domain.Route) model.Route {
	return model.Route{
		RouteID:        d.RouteID,
		RouteShortName: d.RouteShortName,
	}
}

func calendarFromDomain(d domain.Calendar) model.Calendar {
	return model.Calendar{
		ServiceID: d.ServiceID,
		Monday:    d.Monday,
		Tuesday:   d.Tuesday,
		Wednesday: d.Wednesday,
		Thursday:  d.Thursday,
		Friday:    d.Friday,
		Saturday:  d.Saturday,
		Sunday:    d.Sunday,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
	}
}

func calendarDateFromDomain(d domain.CalendarDate) model.CalendarDate {
	return model.CalendarDate{
		ServiceID:     d.ServiceID,
		Date:          d.Date,
		ExceptionType: d.ExceptionType,
	}
}

func tripFromDomain(d domain.Trip) model.Trip {
	return model.Trip{
		TripID:      d.TripID,
		RouteID:     d.RouteID,
		ServiceID:   d.ServiceID,
		DirectionID: d.DirectionID,
	}
}

func stopTimeFromDomain(d domain.StopTime) model.StopTime {
	return model.StopTime{
		TripID:        d.TripID,
		ArrivalTime:   d.ArrivalTime,
		DepartureTime: d.DepartureTime,
		StopID:        d.StopID,
		StopSequence:  d.StopSequence,
	}
}

func fareRuleFromDomain(d domain.FareRule) model.FareRule {
	return model.FareRule{
		RouteID:       d.RouteID,
		OriginID:      d.OriginID,
		DestinationID: d.DestinationID,
		Price:         d.Price,
	}
}

func stopsFromDomain(ds []domain.Stop) []model.Stop {
	out := make([]model.Stop, len(ds))
	for i, d := range ds {
		out[i] = stopFromDomain(d)
	}
	return out
}

func routesFromDomain(ds []domain.Route) []model.Route {
	out := make([]model.Route, len(ds))
	for i, d := range ds {
		out[i] = routeFromDomain(d)
	}
	return out
}

func calendarsFromDomain(ds []domain.Calendar) []model.Calendar {
	out := make([]model.Calendar, len(ds))
	for i, d := range ds {
		out[i] = calendarFromDomain(d)
	}
	return out
}

func calendarDatesFromDomain(ds []domain.CalendarDate) []model.CalendarDate {
	out := make([]model.CalendarDate, len(ds))
	for i, d := range ds {
		out[i] = calendarDateFromDomain(d)
	}
	return out
}

func tripsFromDomain(ds []domain.Trip) []model.Trip {
	out := make([]model.Trip, len(ds))
	for i, d := range ds {
		out[i] = tripFromDomain(d)
	}
	return out
}

func stopTimesFromDomain(ds []domain.StopTime) []model.StopTime {
	out := make([]model.StopTime, len(ds))
	for i, d := range ds {
		out[i] = stopTimeFromDomain(d)
	}
	return out
}

func fareRulesFromDomain(ds []domain.FareRule) []model.FareRule {
	out := make([]model.FareRule, len(ds))
	for i, d := range ds {
		out[i] = fareRuleFromDomain(d)
	}
	return out
}
