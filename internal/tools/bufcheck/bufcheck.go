// Package bufcheck は、buf lint（STANDARD）と Google API Linter では検査できない
// API デザインガイドライン（.rulesync/rules/api-design.md）の規約を検査する buf の check plugin のルールを定義する。
//
// cmd/buf-plugin-dotto から buf.yaml の plugins に登録して `buf lint` で実行する。
package bufcheck

import (
	"buf.build/go/bufplugin/check"
)

// CategoryID は本プラグインのルールをまとめるカテゴリ。buf.yaml の lint.use に指定する。
const CategoryID = "DOTTO"

// Spec は本プラグインの check.Spec。
var Spec = &check.Spec{
	Rules: []*check.RuleSpec{
		updateMaskFieldNumberRuleSpec,
		listPaginationFieldNumbersRuleSpec,
		resourceTimestampsLastRuleSpec,
		fileNoGoPackageRuleSpec,
		fileLayoutRuleSpec,
	},
	Categories: []*check.CategorySpec{
		{
			ID:      CategoryID,
			Purpose: "Checks the dotto API design guidelines that buf lint and Google API Linter do not cover.",
		},
	},
}
