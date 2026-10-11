// Command buf-plugin-dotto は API デザインガイドラインを検査する buf check プラグインである。
//
//	mise run lint:proto  (buf lint のプラグインとして実行する)
package main

import (
	"buf.build/go/bufplugin/check"

	"github.com/fun-dotto/server/internal/tools/protolint"
)

func main() {
	check.Main(protolint.Spec)
}
