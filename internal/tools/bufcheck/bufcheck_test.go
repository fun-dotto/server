package bufcheck_test

import (
	"testing"

	"buf.build/go/bufplugin/check/checktest"

	"github.com/fun-dotto/server/internal/tools/bufcheck"
)

func TestSpec(t *testing.T) {
	checktest.SpecTest(t, bufcheck.Spec)
}

func TestRules(t *testing.T) {
	tests := []struct {
		name   string
		ruleID string
		dir    string
		want   []checktest.ExpectedAnnotation
	}{
		{name: "update_mask が 99", ruleID: bufcheck.UpdateMaskFieldNumberRuleID, dir: "valid"},
		{
			name:   "update_mask が 99 以外",
			ruleID: bufcheck.UpdateMaskFieldNumberRuleID,
			dir:    "update_mask_number",
			want:   []checktest.ExpectedAnnotation{annotation(bufcheck.UpdateMaskFieldNumberRuleID, 12, 2, 12, 44)},
		},
		{name: "ページングのフィールド番号が規約どおり", ruleID: bufcheck.ListPaginationFieldNumbersRuleID, dir: "valid"},
		{name: "ページングしない List", ruleID: bufcheck.ListPaginationFieldNumbersRuleID, dir: "no_paging"},
		{
			name:   "ページングのフィールド番号が規約と異なる",
			ruleID: bufcheck.ListPaginationFieldNumbersRuleID,
			dir:    "pagination_numbers",
			want: []checktest.ExpectedAnnotation{
				annotation(bufcheck.ListPaginationFieldNumbersRuleID, 13, 2, 13, 24),
				annotation(bufcheck.ListPaginationFieldNumbersRuleID, 14, 2, 14, 22),
				annotation(bufcheck.ListPaginationFieldNumbersRuleID, 18, 2, 18, 29),
				annotation(bufcheck.ListPaginationFieldNumbersRuleID, 19, 2, 19, 26),
			},
		},
		{name: "タイムスタンプが末尾", ruleID: bufcheck.ResourceTimestampsLastRuleID, dir: "valid"},
		{
			name:   "タイムスタンプの後にフィールドがある・順序が逆",
			ruleID: bufcheck.ResourceTimestampsLastRuleID,
			dir:    "timestamps_not_last",
			want: []checktest.ExpectedAnnotation{
				annotation(bufcheck.ResourceTimestampsLastRuleID, 8, 2, 8, 43),
				annotation(bufcheck.ResourceTimestampsLastRuleID, 16, 2, 16, 43),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checktest.CheckTest{
				Request: &checktest.RequestSpec{
					Files: &checktest.ProtoFileSpec{
						DirPaths:  []string{"testdata/" + tt.dir},
						FilePaths: []string{"book.proto"},
					},
					RuleIDs: []string{tt.ruleID},
				},
				Spec:                bufcheck.Spec,
				ExpectedAnnotations: tt.want,
			}.Run(t)
		})
	}
}

func annotation(ruleID string, startLine, startColumn, endLine, endColumn int) checktest.ExpectedAnnotation {
	return checktest.ExpectedAnnotation{
		RuleID: ruleID,
		FileLocation: &checktest.ExpectedFileLocation{
			FileName:    "book.proto",
			StartLine:   startLine,
			StartColumn: startColumn,
			EndLine:     endLine,
			EndColumn:   endColumn,
		},
	}
}
