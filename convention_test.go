package convention_test

import (
	"convention"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	result := analysistest.Run(t, testdata, convention.DiagnosticAnalyzer, "a")
	for _, res := range result {
		if diagnose, ok := res.Result.([]analysis.Diagnostic); ok {
			if len(diagnose) != 1 {
				t.Errorf("got len(diagnose) = %d, but want %d", len(diagnose), 1)
			}
			d := diagnose[0]
			pos := res.Pass.Fset.Position(d.Pos)
			if d.Pos != 276 {
				want := res.Pass.Fset.Position(276)
				t.Errorf("got pos = %v, but want %v", pos, want)
			}
			want := "warning: Handle error case or leave comment."
			if d.Message != want {
				t.Errorf("got d.Message = %v, but want %v", d.Message, want)
			}
		}
	}
}
