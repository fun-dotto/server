// Command gormuuidlint は gormuuid アナライザを単体で実行する。
//
//	mise run lint:gormuuid  (go vet -vettool 経由で実行する)
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/fun-dotto/server/internal/tools/gormuuid"
)

func main() {
	singlechecker.Main(gormuuid.Analyzer)
}
