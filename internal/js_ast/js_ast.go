package js_ast

type AST struct {
	Parts []Part
}

type Stmt struct {
	Data S
}

type S interface{ isStmt() }

type Part struct {
	Stmts []Stmt
}

type SEmpty struct{}

func (*SEmpty) isStmt() {}

var SEmptyShared = &SEmpty{}
