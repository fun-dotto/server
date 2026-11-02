package handler

import (
	"fmt"
	"time"
	"uuid"

	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dateLayout = "2006-01-02"

func requireString(v *violations, field, s string) string {
	if s == "" {
		v.add(field, "must not be empty")
	}
	return s
}

func parseUUID(v *violations, field, s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		v.add(field, "must be a UUID")
	}
	return id
}

// parseIDField は Get / Update / Delete で受け取る UUID の識別子を検証する。
func parseIDField(field, s string) (string, error) {
	var v violations
	id := parseUUID(&v, field, s)
	if err := v.err(); err != nil {
		return "", err
	}
	return id.String(), nil
}

func requireTimestamp(v *violations, field string, ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		v.add(field, "must be set")
		return time.Time{}
	}
	if err := ts.CheckValid(); err != nil {
		v.add(field, "must be a valid timestamp")
	}
	return ts.AsTime()
}

func optionalTimestamp(v *violations, field string, ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := requireTimestamp(v, field, ts)
	return &t
}

func toTimestamp(t time.Time) *timestamppb.Timestamp { return timestamppb.New(t) }

func toOptionalTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func requireDate(v *violations, field string, d *date.Date) time.Time {
	if d == nil {
		v.add(field, "must be set")
		return time.Time{}
	}
	t := time.Date(int(d.GetYear()), time.Month(d.GetMonth()), int(d.GetDay()), 0, 0, 0, 0, time.UTC)
	if d.GetYear() < 1 || t.Year() != int(d.GetYear()) || t.Month() != time.Month(d.GetMonth()) || t.Day() != int(d.GetDay()) {
		v.add(field, "must be a valid full date")
	}
	return t
}

// dateKey は Get / Delete で受け取る日付の識別子を検証し、DB に渡す文字列にする。
func dateKey(field string, d *date.Date) (string, error) {
	var v violations
	t := requireDate(&v, field, d)
	if err := v.err(); err != nil {
		return "", err
	}
	return t.Format(dateLayout), nil
}

func toDate(t time.Time) *date.Date {
	return &date.Date{Year: int32(t.Year()), Month: int32(t.Month()), Day: int32(t.Day())}
}

func ptr[T any](v T) *T { return &v }

func int32Ptr(p *int) *int32 {
	if p == nil {
		return nil
	}
	return ptr(int32(*p))
}

func intPtr(p *int32) *int {
	if p == nil {
		return nil
	}
	return ptr(int(*p))
}

// enumMap は proto の enum 値と DB に保存する文字列の対応。
type enumMap[E ~int32] map[E]string

// toDB は enum 値を DB の文字列にする。required の場合は UNSPECIFIED を不正とする。
func (m enumMap[E]) toDB(v *violations, field string, e E, required bool) string {
	if e == 0 && !required {
		return ""
	}
	s, ok := m[e]
	if !ok {
		v.add(field, "must be a valid enum value other than UNSPECIFIED")
	}
	return s
}

func (m enumMap[E]) toOptionalDB(v *violations, field string, e *E) *string {
	if e == nil {
		return nil
	}
	return ptr(m.toDB(v, field, *e, true))
}

func (m enumMap[E]) fromDB(s string) E {
	for e, str := range m {
		if str == s {
			return e
		}
	}
	return 0
}

func (m enumMap[E]) fromOptionalDB(s *string) *E {
	if s == nil {
		return nil
	}
	return ptr(m.fromDB(*s))
}

var gradeEnum = enumMap[adminv1.Grade]{
	adminv1.Grade_GRADE_B1: "B1",
	adminv1.Grade_GRADE_B2: "B2",
	adminv1.Grade_GRADE_B3: "B3",
	adminv1.Grade_GRADE_B4: "B4",
	adminv1.Grade_GRADE_M1: "M1",
	adminv1.Grade_GRADE_M2: "M2",
	adminv1.Grade_GRADE_D1: "D1",
	adminv1.Grade_GRADE_D2: "D2",
	adminv1.Grade_GRADE_D3: "D3",
}

var courseEnum = enumMap[adminv1.Course]{
	adminv1.Course_COURSE_INFORMATION_SYSTEM: "InformationSystem",
	adminv1.Course_COURSE_INFORMATION_DESIGN: "InformationDesign",
	adminv1.Course_COURSE_ADVANCED_ICT:       "AdvancedICT",
	adminv1.Course_COURSE_COMPLEX_SYSTEM:     "ComplexSystem",
	adminv1.Course_COURSE_INTELLIGENT_SYSTEM: "IntelligentSystem",
}

var classEnum = enumMap[adminv1.Class]{
	adminv1.Class_CLASS_A: "A",
	adminv1.Class_CLASS_B: "B",
	adminv1.Class_CLASS_C: "C",
	adminv1.Class_CLASS_D: "D",
	adminv1.Class_CLASS_E: "E",
	adminv1.Class_CLASS_F: "F",
	adminv1.Class_CLASS_G: "G",
	adminv1.Class_CLASS_H: "H",
	adminv1.Class_CLASS_I: "I",
	adminv1.Class_CLASS_J: "J",
	adminv1.Class_CLASS_K: "K",
	adminv1.Class_CLASS_L: "L",
}

var periodEnum = enumMap[adminv1.Period]{
	adminv1.Period_PERIOD_1: "Period1",
	adminv1.Period_PERIOD_2: "Period2",
	adminv1.Period_PERIOD_3: "Period3",
	adminv1.Period_PERIOD_4: "Period4",
	adminv1.Period_PERIOD_5: "Period5",
	adminv1.Period_PERIOD_6: "Period6",
}

var dayOfWeekEnum = enumMap[dayofweek.DayOfWeek]{
	dayofweek.DayOfWeek_MONDAY:    "Monday",
	dayofweek.DayOfWeek_TUESDAY:   "Tuesday",
	dayofweek.DayOfWeek_WEDNESDAY: "Wednesday",
	dayofweek.DayOfWeek_THURSDAY:  "Thursday",
	dayofweek.DayOfWeek_FRIDAY:    "Friday",
	dayofweek.DayOfWeek_SATURDAY:  "Saturday",
	dayofweek.DayOfWeek_SUNDAY:    "Sunday",
}

// requireKeyString は外部由来の識別子（GTFS の ID など）を検証する。
func requireKeyString(field, s string) (string, error) {
	var v violations
	requireString(&v, field, s)
	return s, v.err()
}

func requireNonNegative(v *violations, field string, n int32) int {
	if n < 0 {
		v.add(field, "must not be negative")
	}
	return int(n)
}

func requirePositive(v *violations, field string, n int32) int {
	if n <= 0 {
		v.add(field, "must be positive")
	}
	return int(n)
}

// requireKeyInt は識別子を構成する整数（年度・停車順など）を検証する。
func requireKeyInt(field string, n int32, minimum int32) (int, error) {
	var v violations
	if n < minimum {
		v.add(field, fmt.Sprintf("must be greater than or equal to %d", minimum))
	}
	return int(n), v.err()
}
