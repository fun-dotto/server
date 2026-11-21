// Command buf-plugin-dotto は API デザインガイドラインを検査する buf の check plugin。
// buf.yaml の plugins に登録して `buf lint` から実行する。
package main

import (
	"buf.build/go/bufplugin/check"

	"github.com/fun-dotto/server/internal/tools/bufcheck"
)

func main() {
	check.Main(bufcheck.Spec)
}
