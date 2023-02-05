package nilmiru

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "nilmiru is a static analysis tool that detects nil check leakage in function"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "nilmiru",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	if pass == nil {
		return nil, nil
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	varIsCheckedTable := map[types.Object]bool{}
	inspect.Preorder(nil, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			checkFuncDecl(pass, n, varIsCheckedTable)
		case *ast.Ident:
			checkIdent(pass, n, varIsCheckedTable)
		case *ast.IfStmt:
			checkIfStmt(pass, n, varIsCheckedTable)
		}
	})
	return nil, nil
}

// store arguments which require nil check
func checkFuncDecl(pass *analysis.Pass, funcDecl *ast.FuncDecl, varIsCheckedTable map[types.Object]bool) {
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			obj := pass.TypesInfo.ObjectOf(name)
			switch obj.Type().(type) {
			case *types.Pointer, *types.Slice:
				varIsCheckedTable[pass.TypesInfo.Defs[name]] = false
			}
		}
	}
}

// check if the variable is used after nil check
func checkIdent(pass *analysis.Pass, ident *ast.Ident, varIsCheckedTable map[types.Object]bool) {
	for obj := range varIsCheckedTable {
		if ident_obj := pass.TypesInfo.Uses[ident]; ident_obj != nil {
			if ident_obj == obj {
				if !varIsCheckedTable[obj] {
					pass.Reportf(ident.Pos(), "nil check leakage")
				}
			}
		}
	}
}

// update nil check table
func checkIfStmt(pass *analysis.Pass, ifStmt *ast.IfStmt, varIsCheckedTable map[types.Object]bool) {
	switch ifStmt := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		switch X := ifStmt.X.(type) {
		case *ast.Ident:
			if types.Identical(pass.TypesInfo.TypeOf(ifStmt.Y), types.Typ[types.UntypedNil]) {
				obj := pass.TypesInfo.ObjectOf(X)
				varIsCheckedTable[obj] = true
			}
		case *ast.CallExpr:
			if name, ok := X.Fun.(*ast.Ident); ok {
				if name.Name == "len" {
					if types.Identical(pass.TypesInfo.TypeOf(ifStmt.Y), types.Typ[types.Int]) {
						if ident, ok := X.Args[0].(*ast.Ident); ok {
							obj := pass.TypesInfo.ObjectOf(ident)
							varIsCheckedTable[obj] = true
						}
					}
				}
			}
		}
	}
}
