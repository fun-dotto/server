package model

import (
	"strconv"
	"strings"
)

// ParseSyllabusCredit はシラバスの単位数（例: `2単位`）から数値を取り出す。取り出せない場合は 0。
func ParseSyllabusCredit(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "単位")))
	if err != nil {
		return 0
	}
	return n
}
