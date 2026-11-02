package service

import (
	"slices"
	"testing"
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func TestDetermineSemestersFromDates(t *testing.T) {
	terms := []domain.Term{
		{Year: 2026, Semester: domain.CourseSemesterH1, Start: "2026-04-06", End: "2026-08-03"},
		{Year: 2026, Semester: domain.CourseSemesterQ1, Start: "2026-04-06", End: "2026-06-04"},
		{Year: 2026, Semester: domain.CourseSemesterQ2, Start: "2026-06-05", End: "2026-08-03"},
		{Year: 2026, Semester: domain.CourseSemesterH2, Start: "2026-09-24", End: "2027-01-29"},
		{Year: 2026, Semester: domain.CourseSemesterQ4, Start: "2026-11-24", End: "2027-01-29"},
		{Year: 2027, Semester: domain.CourseSemesterQ1, Start: "2027-04-05", End: "2027-06-03"},
	}
	date := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name          string
		dates         []time.Time
		wantYear      *int
		wantSemesters []domain.CourseSemester
	}{
		{
			name:  "日付が空の場合は何も返さない",
			dates: nil,
		},
		{
			name:          "1つの授業期間に含まれる日付",
			dates:         []time.Time{date(2026, time.April, 20)},
			wantYear:      intPtr(2026),
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH1, domain.CourseSemesterQ1},
		},
		{
			name:          "期間の開始日と終了日は含まれる",
			dates:         []time.Time{date(2026, time.April, 6), date(2026, time.June, 4)},
			wantYear:      intPtr(2026),
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH1, domain.CourseSemesterQ1},
		},
		{
			name:          "複数の授業期間にまたがる日付",
			dates:         []time.Time{date(2026, time.May, 1), date(2026, time.July, 1)},
			wantYear:      intPtr(2026),
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH1, domain.CourseSemesterQ1, domain.CourseSemesterQ2},
		},
		{
			name:          "暦年をまたいでも同じ年度なら年度を返す",
			dates:         []time.Time{date(2026, time.December, 1), date(2027, time.January, 15)},
			wantYear:      intPtr(2026),
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH2, domain.CourseSemesterQ4},
		},
		{
			name:          "複数の年度にまたがる場合は年度を返さない",
			dates:         []time.Time{date(2026, time.May, 1), date(2027, time.April, 20)},
			wantYear:      nil,
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH1, domain.CourseSemesterQ1},
		},
		{
			name:          "授業期間に含まれない日付は無視する",
			dates:         []time.Time{date(2026, time.August, 20), date(2026, time.May, 1)},
			wantYear:      intPtr(2026),
			wantSemesters: []domain.CourseSemester{domain.CourseSemesterH1, domain.CourseSemesterQ1},
		},
		{
			name:  "すべて授業期間外の場合は何も返さない",
			dates: []time.Time{date(2026, time.August, 20)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, semesters := determineSemestersFromDates(terms, tt.dates)

			if (year == nil) != (tt.wantYear == nil) || (year != nil && *year != *tt.wantYear) {
				t.Errorf("year = %v, want %v", year, tt.wantYear)
			}

			slices.Sort(semesters)
			slices.Sort(tt.wantSemesters)
			if !slices.Equal(semesters, tt.wantSemesters) {
				t.Errorf("semesters = %v, want %v", semesters, tt.wantSemesters)
			}
		})
	}
}
