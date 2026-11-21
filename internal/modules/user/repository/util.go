package repository

import (
	"uuid"

	"github.com/fun-dotto/server/internal/modules/user/domain"
)

// canonicalUUIDs は UUID として解釈できる値を正規形の文字列にして返す。
// 解釈できない値は uuid 列との比較でエラーになるため除く（該当する行はない）。
func canonicalUUIDs(ids []string) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		result = append(result, u.String())
	}
	return result
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	result := make([]string, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func uniqueTargetUsers(targets []domain.NotificationTargetUser) []domain.NotificationTargetUser {
	seen := make(map[string]struct{}, len(targets))
	result := make([]domain.NotificationTargetUser, 0, len(targets))
	for _, t := range targets {
		if _, ok := seen[t.UserID]; ok {
			continue
		}
		seen[t.UserID] = struct{}{}
		result = append(result, t)
	}
	return result
}
