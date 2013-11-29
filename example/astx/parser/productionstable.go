package parser

import "code.google.com/p/gocc/example/astx/ast"

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]interface{}) (interface{}, error)
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String:     `S' : StmtList ;`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String:     `StmtList : Stmt << ast.NewStmtList(X[0]) >> ;`,
		Id:         "StmtList",
		NTType:     1,
		Index:      1,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
			return ast.NewStmtList(X[0])
		},
	},
	ProdTabEntry{
		String:     `StmtList : StmtList Stmt << ast.AppendStmt(X[0], X[1]) >> ;`,
		Id:         "StmtList",
		NTType:     1,
		Index:      2,
		NumSymbols: 2,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
			return ast.AppendStmt(X[0], X[1])
		},
	},
	ProdTabEntry{
		String:     `Stmt : id << ast.NewStmt(X[0]) >> ;`,
		Id:         "Stmt",
		NTType:     2,
		Index:      3,
		NumSymbols: 1,
		ReduceFunc: func(X []interface{}) (interface{}, error) {
			return ast.NewStmt(X[0])
		},
	},
}
