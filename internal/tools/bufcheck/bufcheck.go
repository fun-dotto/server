// Package bufcheck は、buf lint（STANDARD）と Google API Linter では検査できない
// API デザインガイドライン（.rulesync/rules/api-design.md）の規約を検査する buf の check plugin のルールを定義する。
package bufcheck

import (
	"context"
	"strings"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// CategoryID は buf.yaml の lint.use に指定する。
const CategoryID = "DOTTO"

const (
	UpdateMaskFieldNumberRuleID      = "DOTTO_UPDATE_MASK_FIELD_NUMBER"
	ListPaginationFieldNumbersRuleID = "DOTTO_LIST_PAGINATION_FIELD_NUMBERS"
	ResourceTimestampsLastRuleID     = "DOTTO_RESOURCE_TIMESTAMPS_LAST"
)

// フィールドの有無と型は Google API Linter（core::0134::request-mask-required・request-mask-field、
// core::0158::*）が検査するため、ここではフィールドがある場合の番号だけを見る。
// App API のようにページングしない List はフィールドを持たないので対象外になる。
var (
	updateMaskFieldNumberRuleSpec = &check.RuleSpec{
		ID:          UpdateMaskFieldNumberRuleID,
		CategoryIDs: []string{CategoryID},
		Purpose:     "Checks that update_mask of Update requests is field number 99.",
		Type:        check.RuleTypeLint,
		Handler:     checkutil.NewMethodRuleHandler(checkUpdateMaskFieldNumber, checkutil.WithoutImports()),
	}
	listPaginationFieldNumbersRuleSpec = &check.RuleSpec{
		ID:          ListPaginationFieldNumbersRuleID,
		CategoryIDs: []string{CategoryID},
		Purpose:     "Checks the field numbers of pagination fields in List requests and responses.",
		Type:        check.RuleTypeLint,
		Handler:     checkutil.NewMethodRuleHandler(checkListPaginationFieldNumbers, checkutil.WithoutImports()),
	}
	resourceTimestampsLastRuleSpec = &check.RuleSpec{
		ID:          ResourceTimestampsLastRuleID,
		CategoryIDs: []string{CategoryID},
		Purpose:     "Checks that created_at and updated_at are declared last in this order.",
		Type:        check.RuleTypeLint,
		Handler:     checkutil.NewMessageRuleHandler(checkResourceTimestampsLast, checkutil.WithoutImports()),
	}
)

var Spec = &check.Spec{
	Rules: []*check.RuleSpec{
		updateMaskFieldNumberRuleSpec,
		listPaginationFieldNumbersRuleSpec,
		resourceTimestampsLastRuleSpec,
	},
	Categories: []*check.CategorySpec{
		{
			ID:      CategoryID,
			Purpose: "Checks the dotto API design guidelines that buf lint and Google API Linter do not cover.",
		},
	},
}

func checkUpdateMaskFieldNumber(
	_ context.Context,
	w check.ResponseWriter,
	_ check.Request,
	method protoreflect.MethodDescriptor,
) error {
	if !strings.HasPrefix(string(method.Name()), "Update") {
		return nil
	}
	checkFieldNumber(w, method.Input(), "update_mask", 99)
	return nil
}

func checkListPaginationFieldNumbers(
	_ context.Context,
	w check.ResponseWriter,
	_ check.Request,
	method protoreflect.MethodDescriptor,
) error {
	if !strings.HasPrefix(string(method.Name()), "List") {
		return nil
	}
	checkFieldNumber(w, method.Input(), "page_size", 1)
	checkFieldNumber(w, method.Input(), "page_token", 2)

	output := method.Output()
	if output.Fields().ByName("next_page_token") == nil {
		return nil
	}
	checkFieldNumber(w, output, "next_page_token", 2)
	fields := output.Fields()
	for i := range fields.Len() {
		field := fields.Get(i)
		if field.IsList() && field.Kind() == protoreflect.MessageKind && field.Number() != 1 {
			w.AddAnnotation(
				check.WithMessagef("Field %q of %q holds the listed resources and must be field number 1.", field.Name(), output.Name()),
				check.WithDescriptor(field),
			)
		}
	}
	return nil
}

func checkFieldNumber(w check.ResponseWriter, message protoreflect.MessageDescriptor, name protoreflect.Name, want protoreflect.FieldNumber) {
	field := message.Fields().ByName(name)
	if field == nil || field.Number() == want {
		return
	}
	w.AddAnnotation(
		check.WithMessagef("Field %q of %q must be field number %d, but is %d.", name, message.Name(), want, field.Number()),
		check.WithDescriptor(field),
	)
}

// フィールド番号ではなく宣言順を見る。既存リソースにフィールドを追加すると番号は updated_at より大きくなるが、
// 番号は変更できないため、宣言を created_at の前に置けば規約を満たすとする。
func checkResourceTimestampsLast(
	_ context.Context,
	w check.ResponseWriter,
	_ check.Request,
	message protoreflect.MessageDescriptor,
) error {
	fields := message.Fields()
	var want []protoreflect.Name
	for _, name := range []protoreflect.Name{"created_at", "updated_at"} {
		if fields.ByName(name) != nil {
			want = append(want, name)
		}
	}
	offset := fields.Len() - len(want)
	for i, name := range want {
		if fields.Get(offset+i).Name() != name {
			w.AddAnnotation(
				check.WithMessagef("Message %q must declare %s at the end of its fields.", message.Name(), strings.Join(namesToStrings(want), " and ")),
				check.WithDescriptor(fields.ByName(name)),
			)
			return nil
		}
	}
	return nil
}

func namesToStrings(names []protoreflect.Name) []string {
	s := make([]string, len(names))
	for i, n := range names {
		s[i] = string(n)
	}
	return s
}
