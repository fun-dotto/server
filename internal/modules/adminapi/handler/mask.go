package handler

import (
	"slices"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// fieldSet は値を反映するフィールド名の集合。nil はすべてのフィールドを表す（Create 用）。
type fieldSet map[string]bool

func (s fieldSet) has(field string) bool { return s == nil || s[field] }

// updateFields は Update で更新するフィールドを決める。
// update_mask が空なら、リクエストで値が入っているフィールドをすべて対象にする。
// 識別子のフィールドと update_mask 自体は対象外とする。
func updateFields(req proto.Message, mask *fieldmaskpb.FieldMask, keyFields ...string) (fieldSet, error) {
	fixed := append(slices.Clone(keyFields), "update_mask")
	msg := req.ProtoReflect()
	fields := msg.Descriptor().Fields()
	set := fieldSet{}
	if len(mask.GetPaths()) == 0 {
		for i := range fields.Len() {
			fd := fields.Get(i)
			if !slices.Contains(fixed, string(fd.Name())) && msg.Has(fd) {
				set[string(fd.Name())] = true
			}
		}
		return set, nil
	}
	var v violations
	for _, p := range mask.GetPaths() {
		fd := fields.ByName(protoreflect.Name(p))
		if fd == nil || slices.Contains(fixed, p) {
			v.add("update_mask", "unknown or immutable field: "+p)
			continue
		}
		set[p] = true
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	return set, nil
}

// columns は fieldSet を DB の列名（または構造体のフィールド名）に変換する。
// overrides にないフィールドは proto のフィールド名をそのまま列名とする。
func (s fieldSet) columns(overrides map[string][]string) []string {
	cols := make([]string, 0, len(s))
	for f := range s {
		if o, ok := overrides[f]; ok {
			cols = append(cols, o...)
			continue
		}
		cols = append(cols, f)
	}
	slices.Sort(cols)
	return cols
}

// present は proto3 optional フィールドに値が設定されているかを返す。
func present(m proto.Message, field string) bool {
	r := m.ProtoReflect()
	fd := r.Descriptor().Fields().ByName(protoreflect.Name(field))
	return fd != nil && r.Has(fd)
}

// optional は proto3 optional フィールドの値を、未設定なら nil として返す。
func optional[T any](m proto.Message, field string, v T) *T {
	if !present(m, field) {
		return nil
	}
	return &v
}
