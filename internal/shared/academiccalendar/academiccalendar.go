// Package academiccalendar は学年暦（振替授業日・休日・開講時期の期間）のデータを提供する。
package academiccalendar

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "embed"
)

// EventsJSON は行事予定（日付 → 行事名のリスト）の JSON。
//
//go:embed events.json
var EventsJSON []byte

// HolidaysJSON は休日（日付 → 休日名）の JSON。
//
//go:embed holidays.json
var HolidaysJSON []byte

// TermsJSON は開講時期ごとの期間の JSON。
//
//go:embed terms.json
var TermsJSON []byte

// Term は開講時期の期間。Semester は "Q1" などの開講時期、Start / End は YYYY-MM-DD。
type Term struct {
	Year     int
	Semester string
	Start    string
	End      string
}

// Contains は日付（YYYY-MM-DD）が期間内かを返す。
func (t Term) Contains(date string) bool {
	return t.Start <= date && date <= t.End
}

// Calendar は学年暦。
type Calendar struct {
	// SubstituteDays は振替授業日（YYYY-MM-DD → 振替先の曜日）。
	SubstituteDays map[string]time.Weekday
	// Holidays は休日（YYYY-MM-DD）。
	Holidays map[string]struct{}
	// Terms は開講時期の期間。
	Terms []Term
}

var japaneseDays = map[string]time.Weekday{
	"月曜": time.Monday,
	"火曜": time.Tuesday,
	"水曜": time.Wednesday,
	"木曜": time.Thursday,
	"金曜": time.Friday,
	"土曜": time.Saturday,
	"日曜": time.Sunday,
}

var validSemesters = map[string]bool{
	"AllYear": true, "H1": true, "H2": true, "Q1": true, "Q2": true, "Q3": true, "Q4": true,
	"SummerIntensive": true, "WinterIntensive": true,
}

// Load は埋め込みの JSON から学年暦を読み込む。
func Load() (Calendar, error) {
	var events map[string][]string
	if err := json.Unmarshal(EventsJSON, &events); err != nil {
		return Calendar{}, fmt.Errorf("failed to parse events JSON: %w", err)
	}
	substitute := make(map[string]time.Weekday)
	for date, names := range events {
		for _, name := range names {
			if !strings.HasSuffix(name, "振替授業") {
				continue
			}
			for prefix, w := range japaneseDays {
				if strings.HasPrefix(name, prefix) {
					substitute[date] = w
					break
				}
			}
		}
	}

	var holidaysRaw map[string]string
	if err := json.Unmarshal(HolidaysJSON, &holidaysRaw); err != nil {
		return Calendar{}, fmt.Errorf("failed to parse holidays JSON: %w", err)
	}
	holidays := make(map[string]struct{}, len(holidaysRaw))
	for date := range holidaysRaw {
		holidays[date] = struct{}{}
	}

	var termsRaw []struct {
		Year     int    `json:"year"`
		Semester string `json:"semester"`
		Start    string `json:"start"`
		End      string `json:"end"`
	}
	if err := json.Unmarshal(TermsJSON, &termsRaw); err != nil {
		return Calendar{}, fmt.Errorf("failed to parse terms JSON: %w", err)
	}
	terms := make([]Term, 0, len(termsRaw))
	for _, r := range termsRaw {
		if !validSemesters[r.Semester] {
			return Calendar{}, fmt.Errorf("invalid term semester %q (year %d)", r.Semester, r.Year)
		}
		start, err := time.Parse(time.DateOnly, r.Start)
		if err != nil {
			return Calendar{}, fmt.Errorf("invalid term start %q: %w", r.Start, err)
		}
		end, err := time.Parse(time.DateOnly, r.End)
		if err != nil {
			return Calendar{}, fmt.Errorf("invalid term end %q: %w", r.End, err)
		}
		if end.Before(start) {
			return Calendar{}, fmt.Errorf("term end %q is before start %q", r.End, r.Start)
		}
		terms = append(terms, Term{Year: r.Year, Semester: r.Semester, Start: r.Start, End: r.End})
	}

	return Calendar{SubstituteDays: substitute, Holidays: holidays, Terms: terms}, nil
}

var jst = time.FixedZone("JST", 9*3600)

// CurrentAcademicYear は now が属する年度を返す。日本時間の 4 月 1 日〜翌年 3 月 31 日を同一年度とする。
func CurrentAcademicYear(now time.Time) int {
	n := now.In(jst)
	if n.Month() >= time.April {
		return n.Year()
	}
	return n.Year() - 1
}
