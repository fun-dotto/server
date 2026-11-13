// Package config はバッチジョブの環境変数を読む。
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

// TargetYear は環境変数 TARGET_YEAR の年度を返す。未設定なら当年度（日本の年度）。
func TargetYear() (int, error) {
	v := os.Getenv("TARGET_YEAR")
	if v == "" {
		return domain.CurrentAcademicYear(time.Now()), nil
	}
	y, err := strconv.Atoi(v)
	if err != nil || y < 2000 {
		return 0, fmt.Errorf("invalid TARGET_YEAR: %q", v)
	}
	return y, nil
}
