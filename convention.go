package convention

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "convention detects code that goes against conventions and prompting the writing of comments for such code."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "convention",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.IfStmt:
			checkIfCond(pass, n.Cond)
		}
	})

	return nil, nil
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

func checkIfCond(pass *analysis.Pass, cond ast.Expr) {
	ast.Inspect(cond, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		binOp, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		// Only check ==
		if binOp.Op != token.EQL {
			return true
		}
		// x or y should be nil
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
		// x or should implements error interface
		errType := types.Universe.Lookup("error").Type()
		if !types.Implements(t, errType.Underlying().(*types.Interface)) {
			return true
		}
		/* This is not working
		// should have comment
		for _, file := range pass.Files {
			cmap := ast.NewCommentMap(pass.Fset, file, file.Comments)
			cmt, ok := cmap[n]
			if ok {
				// found comment
				fmt.Println(cmt)
				return true
			}
		}
		*/
		pass.Reportf(n.Pos(), "Convention is handling error case by if statement. If not you should comment.")
		return true
	})
}
