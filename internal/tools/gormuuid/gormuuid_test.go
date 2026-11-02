package gormuuid_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/fun-dotto/server/internal/tools/gormuuid"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), gormuuid.Analyzer, "a")
}
