package bufcheck

import (
	"context"
	"fmt"
	"path"
	"slices"
	"strings"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"buf.build/go/bufplugin/descriptor"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	FileNoGoPackageRuleID = "DOTTO_FILE_NO_GO_PACKAGE"
	FileLayoutRuleID      = "DOTTO_FILE_LAYOUT"
)

var (
	// go_package は buf.gen.yaml の managed mode で付与する。plugin には managed mode 適用前の定義が渡る。
	fileNoGoPackageRuleSpec = &check.RuleSpec{
		ID:          FileNoGoPackageRuleID,
		CategoryIDs: []string{CategoryID},
		Purpose:     "Checks that go_package is not set because managed mode sets it.",
		Type:        check.RuleTypeLint,
		Handler:     checkutil.NewFileRuleHandler(checkFileNoGoPackage, checkutil.WithoutImports()),
	}
	fileLayoutRuleSpec = &check.RuleSpec{
		ID:          FileLayoutRuleID,
		CategoryIDs: []string{CategoryID},
		Purpose:     "Checks the declaration order: service, supporting messages, resource message, request/response messages, enums.",
		Type:        check.RuleTypeLint,
		Handler:     checkutil.NewFileRuleHandler(checkFileLayout, checkutil.WithoutImports()),
	}
)

// FileOptions の go_package（フィールド番号 11）。FileDescriptorProto の options はフィールド番号 8。
var goPackageSourcePath = protoreflect.SourcePath{8, 11}

func checkFileNoGoPackage(
	_ context.Context,
	w check.ResponseWriter,
	_ check.Request,
	file descriptor.FileDescriptor,
) error {
	if file.FileDescriptorProto().GetOptions().GetGoPackage() == "" {
		return nil
	}
	w.AddAnnotation(
		check.WithMessage("Do not set go_package. It is set by managed mode in buf.gen.yaml."),
		check.WithFileNameAndSourcePath(file.ProtoreflectFileDescriptor().Path(), goPackageSourcePath),
	)
	return nil
}

type section int

const (
	sectionService section = iota
	sectionSupportingMessage
	sectionResourceMessage
	sectionRequestResponse
	sectionEnum
)

// 標準メソッドの Request / Response の並び順。カスタムメソッドの位置はガイドで決めていないため順序を問わない。
var standardMethodPrefixes = []string{"List", "Get", "Create", "Update", "Delete"}

type declaration struct {
	descriptor protoreflect.Descriptor
	line       int
	section    section
	// 標準メソッドの Request / Response のみ standardMethodPrefixes の添字。それ以外は -1。
	methodOrder int
}

func checkFileLayout(
	_ context.Context,
	w check.ResponseWriter,
	_ check.Request,
	file descriptor.FileDescriptor,
) error {
	fd := file.ProtoreflectFileDescriptor()
	// 1 ファイル 1 リソースで、ファイル名がリソース名（subject.proto → Subject）という規約に依存して
	// リソース message を特定する。
	resourceName := toPascalCase(strings.TrimSuffix(path.Base(fd.Path()), ".proto"))
	locations := fd.SourceLocations()

	var decls []declaration
	add := func(d protoreflect.Descriptor, s section, methodOrder int) {
		decls = append(decls, declaration{
			descriptor:  d,
			line:        locations.ByDescriptor(d).StartLine,
			section:     s,
			methodOrder: methodOrder,
		})
	}
	for i := range fd.Services().Len() {
		add(fd.Services().Get(i), sectionService, -1)
	}
	for i := range fd.Enums().Len() {
		add(fd.Enums().Get(i), sectionEnum, -1)
	}
	for i := range fd.Messages().Len() {
		message := fd.Messages().Get(i)
		name := string(message.Name())
		switch {
		case strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response"):
			add(message, sectionRequestResponse, standardMethodOrder(name))
		case name == resourceName:
			add(message, sectionResourceMessage, -1)
		default:
			add(message, sectionSupportingMessage, -1)
		}
	}
	slices.SortFunc(decls, func(a, b declaration) int { return a.line - b.line })

	// 直前の宣言とだけ比べる。最大値と比べると、1 つの宣言の位置ずれでそれ以降のすべてが報告されるため。
	var prevMethod *declaration
	for i := range decls {
		d := &decls[i]
		switch {
		case i > 0 && d.section < decls[i-1].section:
			addLayoutAnnotation(w, d, &decls[i-1])
		case prevMethod != nil && d.methodOrder >= 0 && d.methodOrder < prevMethod.methodOrder:
			addLayoutAnnotation(w, d, prevMethod)
		}
		if d.methodOrder >= 0 {
			prevMethod = d
		}
	}
	return nil
}

func addLayoutAnnotation(w check.ResponseWriter, d, before *declaration) {
	w.AddAnnotation(
		check.WithMessage(fmt.Sprintf(
			"%q must be declared before %q. Declare service, supporting messages, resource message, "+
				"request/response messages (List, Get, Create, Update, Delete), then enums.",
			d.descriptor.Name(), before.descriptor.Name(),
		)),
		check.WithDescriptor(d.descriptor),
	)
}

func standardMethodOrder(messageName string) int {
	for i, prefix := range standardMethodPrefixes {
		if strings.HasPrefix(messageName, prefix) {
			return i
		}
	}
	return -1
}

func toPascalCase(snake string) string {
	var b strings.Builder
	for word := range strings.SplitSeq(snake, "_") {
		if word == "" {
			continue
		}
		b.WriteString(strings.ToUpper(word[:1]) + word[1:])
	}
	return b.String()
}
