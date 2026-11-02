package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type termJSON struct {
	Year     int    `json:"year"`
	Semester string `json:"semester"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

func LoadTerms(termsJSON []byte) ([]domain.Term, error) {
	var raw []termJSON
	if err := json.Unmarshal(termsJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse terms JSON: %w", err)
	}

	result := make([]domain.Term, 0, len(raw))
	for _, r := range raw {
		semester := domain.CourseSemester(r.Semester)
		switch semester {
		case domain.CourseSemesterAllYear,
			domain.CourseSemesterH1, domain.CourseSemesterH2,
			domain.CourseSemesterQ1, domain.CourseSemesterQ2, domain.CourseSemesterQ3, domain.CourseSemesterQ4,
			domain.CourseSemesterSummerIntensive, domain.CourseSemesterWinterIntensive:
		default:
			return nil, fmt.Errorf("invalid term semester %q (year %d)", r.Semester, r.Year)
		}
		start, err := time.Parse(time.DateOnly, r.Start)
		if err != nil {
			return nil, fmt.Errorf("invalid term start %q: %w", r.Start, err)
		}
		end, err := time.Parse(time.DateOnly, r.End)
		if err != nil {
			return nil, fmt.Errorf("invalid term end %q: %w", r.End, err)
		}
		if end.Before(start) {
			return nil, fmt.Errorf("term end %q is before start %q", r.End, r.Start)
		}
		result = append(result, domain.Term{Year: r.Year, Semester: semester, Start: r.Start, End: r.End})
	}
	return result, nil
}
