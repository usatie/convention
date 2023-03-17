package convention

import (
	"golang.org/x/tools/go/analysis"
)

const doc = "convention detects code that goes against conventions and prompting the writing of comments for such code."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name:     "convention",
	Doc:      doc,
	Run:      report,
	Requires: []*analysis.Analyzer{DiagnoseAnalyzer},
}

func report(pass *analysis.Pass) (any, error) {
	result := pass.ResultOf[DiagnoseAnalyzer].([]analysis.Diagnostic)
	for _, d := range result {
		pass.Report(d)
	}
	return nil, nil
}
