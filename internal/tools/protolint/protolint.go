// Package protolint は proto/ 配下の API 定義が API デザインガイドライン
// (.rulesync/rules/api-design.md) に従っているかを検査する buf check プラグインである。
//
// buf lint (STANDARD) で検査できない規約のみを扱う。
// service を 1 つ持つファイルを「リソースファイル」とみなし、
// service 名から "Service" を除いたものをリソース名として扱う。
package protolint

import (
	"context"
	"slices"
	"strings"
	"unicode"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"buf.build/go/bufplugin/descriptor"
	"buf.build/go/bufplugin/info"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const categoryID = "DOTTO"

const (
	ruleNoGoPackage            = "DOTTO_FILE_NO_GO_PACKAGE"
	ruleImportsSorted          = "DOTTO_FILE_IMPORTS_SORTED"
	ruleServiceMatchesFile     = "DOTTO_SERVICE_MATCHES_FILE"
	ruleTopLevelPrefix         = "DOTTO_TOP_LEVEL_RESOURCE_PREFIX"
	ruleDeclarationOrder       = "DOTTO_DECLARATION_ORDER"
	ruleRPCResourceName        = "DOTTO_RPC_RESOURCE_NAME"
	ruleResourceIDFirst        = "DOTTO_RESOURCE_ID_FIRST"
	ruleResourceTimestampsLast = "DOTTO_RESOURCE_TIMESTAMPS_LAST"
	ruleAtFieldTimestamp       = "DOTTO_FIELD_AT_TIMESTAMP"
	ruleIDFieldString          = "DOTTO_FIELD_ID_STRING"
	ruleRepeatedPlural         = "DOTTO_FIELD_REPEATED_PLURAL"
	ruleNoPrepositions         = "DOTTO_FIELD_NO_PREPOSITIONS"
	ruleListPagination         = "DOTTO_LIST_PAGINATION"
	ruleResponseWrapsResource  = "DOTTO_RESPONSE_WRAPS_RESOURCE"
	ruleDeleteResponseEmpty    = "DOTTO_DELETE_RESPONSE_EMPTY"
	ruleIdentifierRequest      = "DOTTO_GET_DELETE_REQUEST_IDENTIFIER"
	ruleCreateRequestNoID      = "DOTTO_CREATE_REQUEST_NO_ID"
	ruleCreateRequestFields    = "DOTTO_CREATE_REQUEST_FIELDS"
	ruleUpdateRequestFields    = "DOTTO_UPDATE_REQUEST_FIELDS"
)

const (
	fullNameTimestamp = "google.protobuf.Timestamp"
	fullNameFieldMask = "google.protobuf.FieldMask"
	updateMaskNumber  = 99
)

// 標準メソッドの動詞。並びはファイル内の Request / Response の記述順でもある。
var standardVerbs = []string{"List", "Get", "Create", "Update", "Delete"}

// pluralize で扱えない複数形。
var irregularPlurals = map[string]string{
	"Syllabus": "Syllabi",
}

// フィールド名に含めない前置詞 (AIP-140)。api-linter の core::0140::prepositions と同じ。
var prepositions = []string{
	"after", "at", "before", "between", "but", "by", "except", "for", "from", "in", "including",
	"into", "of", "over", "since", "to", "toward", "under", "with", "within", "without",
}

// 前置詞を含むが AIP で使われる慣用的なフィールド名 (google.type.DayOfWeek など)。
var allowedPrepositionNames = []string{"day_of_week"}

// Spec は buf check プラグインの仕様である。
var Spec = &check.Spec{
	Info: &info.Spec{
		Documentation: "dotto の API デザインガイドライン (.rulesync/rules/api-design.md) のうち buf lint STANDARD で検査できない規約を検査する。",
	},
	Categories: []*check.CategorySpec{
		{ID: categoryID, Purpose: "Checks the dotto API design guidelines."},
	},
	Rules: []*check.RuleSpec{
		lintRule(ruleNoGoPackage, "Checks that go_package is not set (managed mode sets it).", checkutil.NewFileRuleHandler(checkNoGoPackage, checkutil.WithoutImports())),
		lintRule(ruleImportsSorted, "Checks that imports are sorted alphabetically.", checkutil.NewFileRuleHandler(checkImportsSorted, checkutil.WithoutImports())),
		lintRule(ruleServiceMatchesFile, "Checks that a file has at most one service named <Resource>Service after the file and defines the resource message.", checkutil.NewFileRuleHandler(checkServiceMatchesFile, checkutil.WithoutImports())),
		lintRule(ruleTopLevelPrefix, "Checks that auxiliary messages and enums are prefixed with the resource name.", resourceHandler(checkTopLevelPrefix)),
		lintRule(ruleDeclarationOrder, "Checks that declarations are ordered as auxiliary types, resource, List/Get/Create/Update/Delete requests and responses, service.", resourceHandler(checkDeclarationOrder)),
		lintRule(ruleRPCResourceName, "Checks that RPCs are named <Verb><Resource> (List uses the plural).", resourceHandler(checkRPCResourceName)),
		lintRule(ruleResourceIDFirst, "Checks that the resource message has a string id as field number 1.", resourceHandler(checkResourceIDFirst)),
		lintRule(ruleResourceTimestampsLast, "Checks that created_at and updated_at are the last fields of the resource message.", resourceHandler(checkResourceTimestampsLast)),
		lintRule(ruleAtFieldTimestamp, "Checks that fields named *_at are google.protobuf.Timestamp.", checkutil.NewFieldRuleHandler(checkAtFieldTimestamp, checkutil.WithoutImports())),
		lintRule(ruleIDFieldString, "Checks that fields named id, *_id or *_ids are string.", checkutil.NewFieldRuleHandler(checkIDFieldString, checkutil.WithoutImports())),
		lintRule(ruleRepeatedPlural, "Checks that repeated fields have plural names.", checkutil.NewFieldRuleHandler(checkRepeatedPlural, checkutil.WithoutImports())),
		lintRule(ruleNoPrepositions, "Checks that field names do not include prepositions, except the `_at` suffix of timestamps.", checkutil.NewFieldRuleHandler(checkNoPrepositions, checkutil.WithoutImports())),
		lintRule(ruleListPagination, "Checks List requests and responses for cursor pagination fields.", resourceHandler(checkListPagination)),
		lintRule(ruleResponseWrapsResource, "Checks that Get/Create/Update responses wrap the resource as field 1.", resourceHandler(checkResponseWrapsResource)),
		lintRule(ruleDeleteResponseEmpty, "Checks that Delete responses are empty.", resourceHandler(checkDeleteResponseEmpty)),
		lintRule(ruleIdentifierRequest, "Checks that Get/Delete requests only have the resource identifier fields.", resourceHandler(checkIdentifierRequest)),
		lintRule(ruleCreateRequestNoID, "Checks that Create requests do not have id (assigned by the server).", resourceHandler(checkCreateRequestNoID)),
		lintRule(ruleCreateRequestFields, "Checks that Create requests have resource fields except timestamps.", resourceHandler(checkCreateRequestFields)),
		lintRule(ruleUpdateRequestFields, "Checks that Update requests have the identifier, resource fields and update_mask = 99.", resourceHandler(checkUpdateRequestFields)),
	},
}

func lintRule(id, purpose string, handler check.RuleHandler) *check.RuleSpec {
	return &check.RuleSpec{
		ID:          id,
		CategoryIDs: []string{categoryID},
		Default:     true,
		Purpose:     purpose,
		Type:        check.RuleTypeLint,
		Handler:     handler,
	}
}

// resourceFile はリソースファイルの service とリソース message の組である。
type resourceFile struct {
	file     protoreflect.FileDescriptor
	service  protoreflect.ServiceDescriptor
	name     string
	resource protoreflect.MessageDescriptor
}

// resourceHandler はリソースファイルに対してのみ f を呼ぶ。
func resourceHandler(f func(check.ResponseWriter, *resourceFile)) check.RuleHandler {
	return checkutil.NewFileRuleHandler(
		func(_ context.Context, w check.ResponseWriter, _ check.Request, fd descriptor.FileDescriptor) error {
			if rf := newResourceFile(fd.ProtoreflectFileDescriptor()); rf != nil {
				f(w, rf)
			}
			return nil
		},
		checkutil.WithoutImports(),
	)
}

func newResourceFile(file protoreflect.FileDescriptor) *resourceFile {
	if file.Services().Len() != 1 {
		return nil
	}
	service := file.Services().Get(0)
	name, ok := strings.CutSuffix(string(service.Name()), "Service")
	if !ok {
		return nil
	}
	resource := file.Messages().ByName(protoreflect.Name(name))
	if resource == nil {
		return nil
	}
	return &resourceFile{file: file, service: service, name: name, resource: resource}
}

// method は verb の標準メソッドを返す。
func (rf *resourceFile) method(verb string) protoreflect.MethodDescriptor {
	name := verb + rf.name
	if verb == "List" {
		name = "List" + rf.plural()
	}
	return rf.service.Methods().ByName(protoreflect.Name(name))
}

func (rf *resourceFile) plural() string {
	if p, ok := irregularPlurals[rf.name]; ok {
		return p
	}
	return pluralize(rf.name)
}

// identifierField はリソースの識別子フィールド (フィールド番号 1) を返す。
func (rf *resourceFile) identifierField() protoreflect.FieldDescriptor {
	return rf.resource.Fields().ByNumber(1)
}

func checkNoGoPackage(_ context.Context, w check.ResponseWriter, _ check.Request, fd descriptor.FileDescriptor) error {
	if fd.FileDescriptorProto().GetOptions().GetGoPackage() != "" {
		w.AddAnnotation(
			check.WithFileNameAndSourcePath(fd.ProtoreflectFileDescriptor().Path(), protoreflect.SourcePath{8, 11}),
			check.WithMessage(`go_package must not be set; it is set by managed mode in buf.gen.yaml.`),
		)
	}
	return nil
}

func checkImportsSorted(_ context.Context, w check.ResponseWriter, _ check.Request, fd descriptor.FileDescriptor) error {
	imports := fd.ProtoreflectFileDescriptor().Imports()
	for i := 1; i < imports.Len(); i++ {
		prev, cur := imports.Get(i-1).Path(), imports.Get(i).Path()
		if prev > cur {
			w.AddAnnotation(
				check.WithFileNameAndSourcePath(fd.ProtoreflectFileDescriptor().Path(), protoreflect.SourcePath{3, int32(i)}),
				check.WithMessagef("Import %q must come before %q; imports must be sorted alphabetically.", cur, prev),
			)
		}
	}
	return nil
}

func checkServiceMatchesFile(_ context.Context, w check.ResponseWriter, _ check.Request, fd descriptor.FileDescriptor) error {
	file := fd.ProtoreflectFileDescriptor()
	services := file.Services()
	if services.Len() == 0 {
		return nil
	}
	if services.Len() > 1 {
		w.AddAnnotation(
			check.WithDescriptor(services.Get(1)),
			check.WithMessage("A file must define at most one service (one resource per file)."),
		)
		return nil
	}
	service := services.Get(0)
	base := strings.TrimSuffix(file.Path()[strings.LastIndex(file.Path(), "/")+1:], ".proto")
	resource := snakeToPascal(base)
	if want := resource + "Service"; string(service.Name()) != want {
		w.AddAnnotation(
			check.WithDescriptor(service),
			check.WithMessagef("Service %q must be named %q after the file name %q.", service.Name(), want, base+".proto"),
		)
		return nil
	}
	if file.Messages().ByName(protoreflect.Name(resource)) == nil {
		w.AddAnnotation(
			check.WithDescriptor(service),
			check.WithMessagef("Resource message %q must be defined in the same file as %q.", resource, service.Name()),
		)
	}
	return nil
}

// rpcMessageNames はファイル内の service が Request / Response として使う message 名を返す。
func rpcMessageNames(rf *resourceFile) map[protoreflect.FullName]bool {
	names := map[protoreflect.FullName]bool{}
	methods := rf.service.Methods()
	for i := range methods.Len() {
		names[methods.Get(i).Input().FullName()] = true
		names[methods.Get(i).Output().FullName()] = true
	}
	return names
}

func checkTopLevelPrefix(w check.ResponseWriter, rf *resourceFile) {
	rpcMessages := rpcMessageNames(rf)
	for _, d := range topLevelTypes(rf.file) {
		if rpcMessages[d.FullName()] || d.FullName() == rf.resource.FullName() {
			continue
		}
		if !strings.HasPrefix(string(d.Name()), rf.name) {
			w.AddAnnotation(
				check.WithDescriptor(d),
				check.WithMessagef("%q must be prefixed with the resource name %q (e.g. %q).", d.Name(), rf.name, rf.name+string(d.Name())),
			)
		}
	}
}

func topLevelTypes(file protoreflect.FileDescriptor) []protoreflect.Descriptor {
	var ds []protoreflect.Descriptor
	for i := range file.Enums().Len() {
		ds = append(ds, file.Enums().Get(i))
	}
	for i := range file.Messages().Len() {
		ds = append(ds, file.Messages().Get(i))
	}
	return ds
}

func checkDeclarationOrder(w check.ResponseWriter, rf *resourceFile) {
	// ランク: 0 = 補助型, 1 = リソース, 2.. = RPC の Request / Response, 最後 = service
	rank := map[protoreflect.FullName]int{rf.resource.FullName(): 1}
	label := map[int]string{0: "auxiliary types", 1: "the resource message"}
	methods := rf.service.Methods()
	var ordered []protoreflect.MethodDescriptor
	for _, verb := range standardVerbs {
		if m := rf.method(verb); m != nil {
			ordered = append(ordered, m)
		}
	}
	for i := range methods.Len() {
		if m := methods.Get(i); !slices.Contains(ordered, m) {
			ordered = append(ordered, m)
		}
	}
	for i, m := range ordered {
		r := 2 + i
		label[r] = "messages of " + string(m.Name())
		for _, msg := range []protoreflect.MessageDescriptor{m.Input(), m.Output()} {
			if msg.ParentFile().Path() == rf.file.Path() {
				if _, ok := rank[msg.FullName()]; !ok {
					rank[msg.FullName()] = r
				}
			}
		}
	}
	serviceRank := 2 + len(ordered)
	label[serviceRank] = "the service"

	type decl struct {
		d    protoreflect.Descriptor
		line int
		rank int
	}
	var decls []decl
	locs := rf.file.SourceLocations()
	for _, d := range topLevelTypes(rf.file) {
		decls = append(decls, decl{d, locs.ByDescriptor(d).StartLine, rank[d.FullName()]})
	}
	decls = append(decls, decl{rf.service, locs.ByDescriptor(rf.service).StartLine, serviceRank})
	slices.SortFunc(decls, func(a, b decl) int { return a.line - b.line })

	maxRank := 0
	for _, d := range decls {
		if d.rank < maxRank {
			w.AddAnnotation(
				check.WithDescriptor(d.d),
				check.WithMessagef("%q (%s) must be declared before %s; order is auxiliary types, resource, List/Get/Create/Update/Delete requests and responses, service.", d.d.Name(), label[d.rank], label[maxRank]),
			)
			continue
		}
		maxRank = d.rank
	}
}

func checkRPCResourceName(w check.ResponseWriter, rf *resourceFile) {
	methods := rf.service.Methods()
	for i := range methods.Len() {
		m := methods.Get(i)
		name := string(m.Name())
		if rest, ok := strings.CutPrefix(name, "List"); ok {
			if rest != rf.plural() {
				w.AddAnnotation(
					check.WithDescriptor(m),
					check.WithMessagef("List RPC %q must be named %q.", name, "List"+rf.plural()),
				)
			}
			continue
		}
		verb, ok := strings.CutSuffix(name, rf.name)
		if !ok || verb == "" {
			w.AddAnnotation(
				check.WithDescriptor(m),
				check.WithMessagef("RPC %q must be named <Verb>%s.", name, rf.name),
			)
		}
	}
}

func checkResourceIDFirst(w check.ResponseWriter, rf *resourceFile) {
	f := rf.identifierField()
	if f == nil || f.Name() != "id" || f.Kind() != protoreflect.StringKind || f.Cardinality() == protoreflect.Repeated {
		w.AddAnnotation(
			check.WithDescriptor(rf.resource),
			check.WithMessagef("Resource %q must have `string id = 1;`.", rf.name),
		)
	}
}

func checkResourceTimestampsLast(w check.ResponseWriter, rf *resourceFile) {
	fields := rf.resource.Fields()
	var want []protoreflect.Name
	for _, name := range []protoreflect.Name{"created_at", "updated_at"} {
		if fields.ByName(name) != nil {
			want = append(want, name)
		}
	}
	if len(want) == 0 {
		return
	}
	tail := fields.Len() - len(want)
	for i, name := range want {
		f := fields.Get(tail + i)
		if f.Name() != name || hasLargerNumber(fields, f) {
			w.AddAnnotation(
				check.WithDescriptor(fields.ByName(name)),
				check.WithMessagef("%q must be placed at the end of %q (created_at, then updated_at) with the largest field numbers.", name, rf.name),
			)
		}
	}
}

func hasLargerNumber(fields protoreflect.FieldDescriptors, f protoreflect.FieldDescriptor) bool {
	for i := range fields.Len() {
		g := fields.Get(i)
		if g.Index() < f.Index() && g.Number() > f.Number() {
			return true
		}
	}
	return false
}

func checkAtFieldTimestamp(_ context.Context, w check.ResponseWriter, _ check.Request, f protoreflect.FieldDescriptor) error {
	if !strings.HasSuffix(string(f.Name()), "_at") {
		return nil
	}
	if f.Message() == nil || f.Message().FullName() != fullNameTimestamp {
		w.AddAnnotation(
			check.WithDescriptor(f),
			check.WithMessagef("Field %q must be %s.", f.Name(), fullNameTimestamp),
		)
	}
	return nil
}

func checkIDFieldString(_ context.Context, w check.ResponseWriter, _ check.Request, f protoreflect.FieldDescriptor) error {
	name := string(f.Name())
	if name != "id" && !strings.HasSuffix(name, "_id") && !strings.HasSuffix(name, "_ids") {
		return nil
	}
	if f.Kind() != protoreflect.StringKind || f.IsMap() {
		w.AddAnnotation(
			check.WithDescriptor(f),
			check.WithMessagef("ID field %q must be string.", f.Name()),
		)
	}
	return nil
}

func checkRepeatedPlural(_ context.Context, w check.ResponseWriter, _ check.Request, f protoreflect.FieldDescriptor) error {
	if f.Cardinality() != protoreflect.Repeated || f.IsMap() {
		return nil
	}
	name := string(f.Name())
	last := name[strings.LastIndex(name, "_")+1:]
	if !strings.HasSuffix(last, "s") && !isIrregularPlural(last) {
		w.AddAnnotation(
			check.WithDescriptor(f),
			check.WithMessagef("Repeated field %q must have a plural name.", f.Name()),
		)
	}
	return nil
}

// checkNoPrepositions は AIP-140 の前置詞の検査のうち、タイムスタンプの `<verb>_at` だけを許可する。
func checkNoPrepositions(_ context.Context, w check.ResponseWriter, _ check.Request, f protoreflect.FieldDescriptor) error {
	if slices.Contains(allowedPrepositionNames, string(f.Name())) {
		return nil
	}
	words := strings.Split(string(f.Name()), "_")
	for i, word := range words {
		if word == "at" && i == len(words)-1 && i > 0 {
			continue
		}
		if slices.Contains(prepositions, word) {
			w.AddAnnotation(
				check.WithDescriptor(f),
				check.WithMessagef("Field %q must not include the preposition %q.", f.Name(), word),
			)
		}
	}
	return nil
}

// isIrregularPlural は snake_case の単語が不規則な複数形かを返す。
func isIrregularPlural(word string) bool {
	for _, p := range irregularPlurals {
		if pascalToSnake(p) == word {
			return true
		}
	}
	return false
}

func checkListPagination(w check.ResponseWriter, rf *resourceFile) {
	m := rf.method("List")
	if m == nil {
		return
	}
	req, resp := m.Input(), m.Output()
	expectField(w, req, "page_size", 1, protoreflect.Int32Kind, "", false)
	expectField(w, req, "page_token", 2, protoreflect.StringKind, "", false)
	expectField(w, resp, protoreflect.Name(pascalToSnake(rf.plural())), 1, protoreflect.MessageKind, rf.resource.FullName(), true)
	expectField(w, resp, "next_page_token", 2, protoreflect.StringKind, "", false)
}

// expectField は msg が `name = number` のフィールドを指定の型で持つことを検査する。
func expectField(w check.ResponseWriter, msg protoreflect.MessageDescriptor, name protoreflect.Name, number protoreflect.FieldNumber, kind protoreflect.Kind, messageName protoreflect.FullName, repeated bool) {
	f := msg.Fields().ByName(name)
	ok := f != nil && f.Number() == number && f.Kind() == kind && (f.Cardinality() == protoreflect.Repeated) == repeated
	if ok && messageName != "" {
		ok = f.Message().FullName() == messageName
	}
	if ok {
		return
	}
	typ := kind.String()
	if messageName != "" {
		typ = string(messageName)
	}
	if repeated {
		typ = "repeated " + typ
	}
	w.AddAnnotation(
		check.WithDescriptor(msg),
		check.WithMessagef("%q must have `%s %s = %d;`.", msg.Name(), typ, name, number),
	)
}

func checkResponseWrapsResource(w check.ResponseWriter, rf *resourceFile) {
	for _, verb := range []string{"Get", "Create", "Update"} {
		m := rf.method(verb)
		if m == nil {
			continue
		}
		resp := m.Output()
		f := resp.Fields().ByNumber(1)
		if resp.Fields().Len() != 1 || f == nil || f.Name() != protoreflect.Name(pascalToSnake(rf.name)) ||
			f.Message() == nil || f.Message().FullName() != rf.resource.FullName() || f.Cardinality() == protoreflect.Repeated {
			w.AddAnnotation(
				check.WithDescriptor(resp),
				check.WithMessagef("%q must only have `%s %s = 1;`.", resp.Name(), rf.name, pascalToSnake(rf.name)),
			)
		}
	}
}

func checkDeleteResponseEmpty(w check.ResponseWriter, rf *resourceFile) {
	m := rf.method("Delete")
	if m == nil {
		return
	}
	if resp := m.Output(); resp.Fields().Len() != 0 {
		w.AddAnnotation(
			check.WithDescriptor(resp),
			check.WithMessagef("%q must be an empty message.", resp.Name()),
		)
	}
}

func checkIdentifierRequest(w check.ResponseWriter, rf *resourceFile) {
	id := rf.identifierField()
	if id == nil {
		return
	}
	for _, verb := range []string{"Get", "Delete"} {
		m := rf.method(verb)
		if m == nil {
			continue
		}
		req := m.Input()
		if req.Fields().ByName(id.Name()) == nil {
			w.AddAnnotation(
				check.WithDescriptor(req),
				check.WithMessagef("%q must have the identifier field %q of %q.", req.Name(), id.Name(), rf.name),
			)
		}
		fields := req.Fields()
		for i := range fields.Len() {
			if f := fields.Get(i); !sameAsResourceField(rf, f) {
				w.AddAnnotation(
					check.WithDescriptor(f),
					check.WithMessagef("%q must only have identifier fields of %q with the same types; %q is not one of them.", req.Name(), rf.name, f.Name()),
				)
			}
		}
	}
}

func checkCreateRequestNoID(w check.ResponseWriter, rf *resourceFile) {
	m := rf.method("Create")
	if m == nil {
		return
	}
	if f := m.Input().Fields().ByName("id"); f != nil {
		w.AddAnnotation(
			check.WithDescriptor(f),
			check.WithMessagef("%q must not have %q; it is assigned by the server.", m.Input().Name(), f.Name()),
		)
	}
}

func checkCreateRequestFields(w check.ResponseWriter, rf *resourceFile) {
	m := rf.method("Create")
	if m == nil {
		return
	}
	fields := m.Input().Fields()
	for i := range fields.Len() {
		f := fields.Get(i)
		switch {
		case f.Name() == "id":
			// DOTTO_CREATE_REQUEST_NO_ID で検査する
		case isTimestampField(f):
			w.AddAnnotation(
				check.WithDescriptor(f),
				check.WithMessagef("%q must not have %q; it is output only.", m.Input().Name(), f.Name()),
			)
		case !sameAsResourceField(rf, f):
			w.AddAnnotation(
				check.WithDescriptor(f),
				check.WithMessagef("%q must be a field of %q with the same type.", f.Name(), rf.name),
			)
		}
	}
}

func checkUpdateRequestFields(w check.ResponseWriter, rf *resourceFile) {
	m := rf.method("Update")
	if m == nil {
		return
	}
	req := m.Input()
	if id := rf.identifierField(); id != nil && req.Fields().ByName(id.Name()) == nil {
		w.AddAnnotation(
			check.WithDescriptor(req),
			check.WithMessagef("%q must have the identifier field %q of %q.", req.Name(), id.Name(), rf.name),
		)
	}
	fields := req.Fields()
	hasMask := false
	for i := range fields.Len() {
		f := fields.Get(i)
		switch {
		case f.Name() == "update_mask":
			hasMask = f.Number() == updateMaskNumber && f.Message() != nil && f.Message().FullName() == fullNameFieldMask
		case isTimestampField(f):
			w.AddAnnotation(
				check.WithDescriptor(f),
				check.WithMessagef("%q must not have %q; it is output only.", req.Name(), f.Name()),
			)
		case !sameAsResourceField(rf, f):
			w.AddAnnotation(
				check.WithDescriptor(f),
				check.WithMessagef("%q must be a field of %q with the same type.", f.Name(), rf.name),
			)
		}
	}
	if !hasMask {
		w.AddAnnotation(
			check.WithDescriptor(req),
			check.WithMessagef("%q must have `%s update_mask = %d;`.", req.Name(), fullNameFieldMask, updateMaskNumber),
		)
	}
}

func isTimestampField(f protoreflect.FieldDescriptor) bool {
	return f.Name() == "created_at" || f.Name() == "updated_at"
}

// sameAsResourceField は f と同名・同型のフィールドがリソースにあるかを返す。
func sameAsResourceField(rf *resourceFile, f protoreflect.FieldDescriptor) bool {
	g := rf.resource.Fields().ByName(f.Name())
	if g == nil || g.Kind() != f.Kind() || g.Cardinality() != f.Cardinality() {
		return false
	}
	if f.Message() != nil && g.Message().FullName() != f.Message().FullName() {
		return false
	}
	if f.Enum() != nil && g.Enum().FullName() != f.Enum().FullName() {
		return false
	}
	return true
}

func pluralize(s string) string {
	switch {
	case strings.HasSuffix(s, "y") && len(s) > 1 && !strings.ContainsRune("aeiou", rune(s[len(s)-2])):
		return s[:len(s)-1] + "ies"
	case strings.HasSuffix(s, "s"), strings.HasSuffix(s, "x"), strings.HasSuffix(s, "ch"), strings.HasSuffix(s, "sh"):
		return s + "es"
	default:
		return s + "s"
	}
}

func snakeToPascal(s string) string {
	var b strings.Builder
	for part := range strings.SplitSeq(s, "_") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	return b.String()
}

func pascalToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}
