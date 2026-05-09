package domain

import (
	"testing"
	"time"
)

func TestCurrentAcademicYear(t *testing.T) {
	jst := time.FixedZone("JST", 9*3600)

	tests := []struct {
		name string
		now  time.Time
		want int
	}{
		{
			name: "JST 1月1日は前暦年の年度",
			now:  time.Date(2026, time.January, 1, 0, 0, 0, 0, jst),
			want: 2025,
		},
		{
			name: "JST 3月31日は前暦年の年度",
			now:  time.Date(2026, time.March, 31, 23, 59, 59, 0, jst),
			want: 2025,
		},
		{
			name: "JST 4月1日から新年度",
			now:  time.Date(2026, time.April, 1, 0, 0, 0, 0, jst),
			want: 2026,
		},
		{
			name: "JST 12月31日は当暦年の年度",
			now:  time.Date(2026, time.December, 31, 23, 59, 59, 0, jst),
			want: 2026,
		},
		{
			name: "UTC 3月31日15時はJSTでは4月1日0時なので新年度",
			now:  time.Date(2026, time.March, 31, 15, 0, 0, 0, time.UTC),
			want: 2026,
		},
		{
			name: "UTC 3月31日14時59分はJSTでは3月31日23時59分なので前年度",
			now:  time.Date(2026, time.March, 31, 14, 59, 59, 0, time.UTC),
			want: 2025,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentAcademicYear(tt.now); got != tt.want {
				t.Errorf("CurrentAcademicYear(%s) = %d, want %d", tt.now.Format(time.RFC3339), got, tt.want)
			}
		})
	}
}
