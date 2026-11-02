// Package gormuuid は GORM のクエリ引数に uuid.UUID を直接渡している箇所を検出する。
//
// 標準 uuid.UUID は [16]byte で driver.Valuer を実装していないため、
// GORM は "(" 直後の ? に渡されると配列として (b0,b1,...,b15) に展開してしまう。
// 位置に依存した不具合を避けるため、GORM に渡す UUID は常に .String() で文字列化する。
package gormuuid

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gormuuid",
	Doc:      "GORM のクエリ引数に uuid.UUID を直接渡している箇所を検出する (.String() で文字列化すること)",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || !isGormDB(pass.TypesInfo.TypeOf(sel.X)) {
			return
		}
		sig, ok := pass.TypesInfo.TypeOf(call.Fun).(*types.Signature)
		if !ok {
			return
		}
		params := sig.Params()
		for i, arg := range call.Args {
			if !isAnyParam(sig, params, i, call.Ellipsis.IsValid()) {
				continue
			}
			if containsUUID(pass.TypesInfo.TypeOf(arg)) {
				pass.Reportf(arg.Pos(), "GORM に uuid.UUID を直接渡さず .String() で文字列化してください")
			}
		}
	})
	return nil, nil
}

// isGormDB は t が gorm.io/gorm.DB またはそのポインタかを返す。
func isGormDB(t types.Type) bool {
	if t == nil {
		return false
	}
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "gorm.io/gorm" && obj.Name() == "DB"
}

// isAnyParam は i 番目の引数が interface{} 型のパラメータに渡されるかを返す。
func isAnyParam(sig *types.Signature, params *types.Tuple, i int, spread bool) bool {
	var t types.Type
	switch {
	case sig.Variadic() && i >= params.Len()-1:
		if spread {
			return false
		}
		t = params.At(params.Len() - 1).Type().(*types.Slice).Elem()
	case i < params.Len():
		t = params.At(i).Type()
	default:
		return false
	}
	iface, ok := t.Underlying().(*types.Interface)
	return ok && iface.Empty()
}

// containsUUID は t が uuid.UUID、またはその slice / array / pointer かを返す。
func containsUUID(t types.Type) bool {
	for {
		switch u := t.(type) {
		case *types.Pointer:
			t = u.Elem()
			continue
		case *types.Slice:
			t = u.Elem()
			continue
		}
		break
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "uuid" && obj.Name() == "UUID"
}
