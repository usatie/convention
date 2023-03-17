package convention

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
)

// Add more diagnosers here
var diagnosers = []*analysis.Analyzer{
	ErrorHandlingDiagnoser,
}

var DiagnoseAnalyzer = &analysis.Analyzer{
	Name:       "diagnose",
	Doc:        "diagnose only. It doesn't report.",
	Run:        diagnose,
	ResultType: reflect.TypeOf([]analysis.Diagnostic{}),
	Requires:   diagnosers,
}

func diagnose(pass *analysis.Pass) (any, error) {
	var combined []analysis.Diagnostic
	for _, diagnoser := range diagnosers {
		if result, ok := pass.ResultOf[diagnoser].([]analysis.Diagnostic); ok {
			combined = append(combined, result...)
		}
	}
	return combined, nil
}
