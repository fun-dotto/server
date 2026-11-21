package handler

import "strings"

// flagSubjectSortByUserAttributes は科目一覧・時間割一覧をユーザー属性順にソートするかの AB テスト用フラグ。
const flagSubjectSortByUserAttributes = "subject_sort_by_user_attributes"

// featureFlagEnabled は X-Flags ヘッダ（`name=true,name2=false` 形式）で name が true のときに true を返す。
// ヘッダ未指定・フラグ未指定・パースできない値はすべて false として扱う。
func featureFlagEnabled(header *string, name string) bool {
	if header == nil {
		return false
	}
	for pair := range strings.SplitSeq(*header, ",") {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || strings.TrimSpace(key) != name {
			continue
		}
		return strings.EqualFold(strings.TrimSpace(value), "true")
	}
	return false
}
