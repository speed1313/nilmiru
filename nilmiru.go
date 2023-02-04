package nilmiru

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "nilmiru is a golang linter which checks nil check leakage in function."

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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	pointerVarList := []types.Object{}
	varIsCheckedTable := map[types.Object]bool{}
	inspect.Preorder(nil, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			pointerVarList = checkFuncDecl(pass, n, pointerVarList)
		case *ast.Ident:
			checkIdent(pass, n, pointerVarList, varIsCheckedTable)

		case *ast.IfStmt:
			checkIfStmt(pass, n, pointerVarList, varIsCheckedTable)
		}

	})

	return nil, nil
}

// store arguments which require nil check
func checkFuncDecl(pass *analysis.Pass, funcDecl *ast.FuncDecl, pointerVarList []types.Object) []types.Object {
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			obj := pass.TypesInfo.ObjectOf(name)
			switch obj.Type().(type) {
			case *types.Pointer:
				pointerVarList = append(pointerVarList, pass.TypesInfo.Defs[name])
			case *types.Slice:
				pointerVarList = append(pointerVarList, pass.TypesInfo.Defs[name])
			}
		}
	}
	return pointerVarList
}

// check if the variable is used after nil check
func checkIdent(pass *analysis.Pass, ident *ast.Ident, pointerVarList []types.Object, varIsCheckedTable map[types.Object]bool) {
	for _, obj := range pointerVarList {
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
func checkIfStmt(pass *analysis.Pass, ifStmt *ast.IfStmt, pointerVarList []types.Object, varIsCheckedTable map[types.Object]bool) {
	switch ifStmt := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		switch ifStmt.X.(type) {
		case *ast.Ident:
			if types.Identical(pass.TypesInfo.TypeOf(ifStmt.Y), types.Typ[types.UntypedNil]) {
				obj := pass.TypesInfo.ObjectOf(ifStmt.X.(*ast.Ident))
				varIsCheckedTable[obj] = true
			}
		case *ast.CallExpr:
			if ifStmt.X.(*ast.CallExpr).Fun.(*ast.Ident).Name == "len" {
				if types.Identical(pass.TypesInfo.TypeOf(ifStmt.Y), types.Typ[types.Int]) {
					obj := pass.TypesInfo.ObjectOf(ifStmt.X.(*ast.CallExpr).Args[0].(*ast.Ident))
					varIsCheckedTable[obj] = true
				}
			}

		}
	}
}
