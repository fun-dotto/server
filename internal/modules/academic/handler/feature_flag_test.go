package handler

import "testing"

func TestFeatureFlagEnabled(t *testing.T) {
	str := func(s string) *string { return &s }
	tests := []struct {
		name   string
		header *string
		want   bool
	}{
		{name: "ヘッダなし", header: nil, want: false},
		{name: "true", header: str("subject_sort_by_user_attributes=true"), want: true},
		{name: "false", header: str("subject_sort_by_user_attributes=false"), want: false},
		{name: "複数フラグ・空白あり", header: str("foo=false, subject_sort_by_user_attributes=true"), want: true},
		{name: "別フラグのみ", header: str("foo=true"), want: false},
		{name: "不正な値", header: str("subject_sort_by_user_attributes"), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := featureFlagEnabled(tt.header, flagSubjectSortByUserAttributes); got != tt.want {
				t.Errorf("featureFlagEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
