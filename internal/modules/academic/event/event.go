package event

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

var japaneseDayToEnum = map[string]domain.DayOfWeek{
	"月曜": domain.DayOfWeekMonday,
	"火曜": domain.DayOfWeekTuesday,
	"水曜": domain.DayOfWeekWednesday,
	"木曜": domain.DayOfWeekThursday,
	"金曜": domain.DayOfWeekFriday,
	"土曜": domain.DayOfWeekSaturday,
	"日曜": domain.DayOfWeekSunday,
}

// LoadSubstituteDayMap はイベント用 JSON を解析し、日付文字列（例: "2026-04-30"）をキーに、
// 「〇曜振替授業」エントリに対応する振替先の DayOfWeek を表すマップを返す。
func LoadSubstituteDayMap(eventsJSON []byte) (map[string]domain.DayOfWeek, error) {
	var raw map[string][]string
	if err := json.Unmarshal(eventsJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse events JSON: %w", err)
	}

	result := make(map[string]domain.DayOfWeek)
	for date, events := range raw {
		for _, e := range events {
			if !strings.HasSuffix(e, "振替授業") {
				continue
			}
			for prefix, dow := range japaneseDayToEnum {
				if strings.HasPrefix(e, prefix) {
					result[date] = dow
					break
				}
			}
		}
	}
	return result, nil
}

// LoadHolidaySet は休日用 JSON を解析し、休日の日付文字列（例: "2026-04-30"）のセットを返す。
func LoadHolidaySet(holidaysJSON []byte) (map[string]struct{}, error) {
	var raw map[string]string
	if err := json.Unmarshal(holidaysJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse holidays JSON: %w", err)
	}

	result := make(map[string]struct{}, len(raw))
	for date := range raw {
		result[date] = struct{}{}
	}
	return result, nil
}
