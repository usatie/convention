package convention

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const doc = "convention detects code that goes against conventions and prompting the writing of comments for such code."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "convention",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var result []analysis.Diagnostic
	for _, file := range pass.Files {
		cmap := ast.NewCommentMap(pass.Fset, file, file.Comments)
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			if n, ok := n.(*ast.IfStmt); ok {
				d, found := checkIfCond(pass, n, cmap)
				if found {
					result = append(result, d)
				}
			}
			return true
		})
	}

	return result, nil
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

func checkIfCond(pass *analysis.Pass, ifStmt *ast.IfStmt, cmap ast.CommentMap) (d analysis.Diagnostic, found bool) {
	ast.Inspect(ifStmt.Cond, func(n ast.Node) bool {
		if n == nil {
			return false
		}
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
		if t == nil {
			return true
		}
		// x or y implements error interface
		errType := types.Universe.Lookup("error").Type()
		if !types.Implements(t, errType.Underlying().(*types.Interface)) {
			return true
		}
		// There is no comments about breaking convention
		for _, cg := range cmap {
			for _, comment := range cg {
				if pass.Fset.Position(comment.Pos()).Line == pass.Fset.Position(ifStmt.Pos()).Line {
					return true
				}
			}
		}
		found = true
		d = analysis.Diagnostic{Pos: ifStmt.Pos(), Message: "warning: Handle error case or leave comment."}
		fmt.Println(pass.Fset.Position(ifStmt.Pos()), ": warning: Handle error case or leave comment.")
		return true
	})
	return d, found
}
