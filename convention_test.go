package convention_test

import (
	"convention"
	"flag"
	"fmt"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/tenntenn/golden"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	result := analysistest.Run(t, testdata, convention.DiagnoseAnalyzer, "a")
	got := ""
	for _, res := range result {
		if diagnose, ok := res.Result.([]analysis.Diagnostic); ok {
			for _, d := range diagnose {
				pos := res.Pass.Fset.Position(d.Pos)
				got += fmt.Sprintf("%s: %s\n", pos, d.Message)
			}
		}
	}
	if diff := golden.Check(t, flagUpdate, "testdata", "a", got); diff != "" {
		t.Error(diff)
	}
}
