package convention

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "convention detects code that goes against conventions and prompting the writing of comments for such code."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name:     "convention",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{DiagnosticAnalyzer},
}

var DiagnosticAnalyzer = &analysis.Analyzer{
	Name:       "diagnostic",
	Doc:        "diagnose only",
	Run:        diagnose,
	ResultType: reflect.TypeOf([]analysis.Diagnostic{}),
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	result := pass.ResultOf[DiagnosticAnalyzer].([]analysis.Diagnostic)
	for _, d := range result {
		pass.Report(d)
	}
	return nil, nil
}

func diagnose(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}
	var result []analysis.Diagnostic
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		if n, ok := n.(*ast.IfStmt); ok {
			d, found := checkIfCond(pass, n)
			if found {
				result = append(result, d)
			}
		}
	})
	return result, nil
}

func checkIfCond(pass *analysis.Pass, ifStmt *ast.IfStmt) (d analysis.Diagnostic, found bool) {
	ast.Inspect(ifStmt.Cond, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		// n is binary expression
		binOp, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		// operator is ==
		if binOp.Op != token.EQL {
			return true
		}
		// x or y is nil
		var t types.Type
		if isNil(pass, binOp.X) {
			t = pass.TypesInfo.TypeOf(binOp.Y)
		} else if isNil(pass, binOp.Y) {
			t = pass.TypesInfo.TypeOf(binOp.X)
		} else {
			return true
		}
		// x or y implements error interface
		if !implementsErrorInterface(t) {
			return true
		}
		// There is no comments about breaking convention
		ifPos := pass.Fset.Position(ifStmt.Pos())
		for _, f := range pass.Files {
			for _, cg := range f.Comments {
				for _, comment := range cg.List {
					if pass.Fset.Position(comment.Pos()).Line == ifPos.Line {
						return true
					}
				}
			}

		}
		found = true
		d = analysis.Diagnostic{
			Pos:     ifStmt.Pos(),
			Message: fmt.Sprintf("warning: Handle error case or leave comment: %s", getFirstLine(ifStmt)),
		}
		return true
	})
	return d, found
}

func isNil(pass *analysis.Pass, n ast.Node) bool {
	ident, ok := n.(*ast.Ident)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.ObjectOf(ident)
	_, ok = obj.(*types.Nil)
	return ok
}

func implementsErrorInterface(t types.Type) bool {
	if t == nil {
		return false
	}
	errType := types.Universe.Lookup("error").Type()
	return types.Implements(t, errType.Underlying().(*types.Interface))
}

func getFirstLine(n ast.Node) string {
	var buf bytes.Buffer
	format.Node(&buf, token.NewFileSet(), n)
	lines := strings.Split(buf.String(), "\n")
	if len(lines) == 0 {
		return ""
	}
	return lines[0]
}
