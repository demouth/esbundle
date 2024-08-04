package js_printer

import (
	"fmt"

	"github.com/demouth/esbundle/internal/js_ast"
)

func Print(tree js_ast.AST, options Options) PrintResult {
	p := printer{
		options: options,
	}
	for _, part := range tree.Parts {
		for _, stmt := range part.Stmts {
			p.printStmt(stmt)
		}
	}
	result := PrintResult{
		JS: p.js,
	}
	return result
}

type Options struct{}

type PrintResult struct {
	JS []byte
}

type printer struct {
	js      []byte
	options Options
}

func (p *printer) printStmt(stmt js_ast.Stmt) {
	switch s := stmt.Data.(type) {
	case *js_ast.SEmpty:
		p.print(";")
		p.printNewline()
	default:
		panic(fmt.Sprintf("Unexpected statement of type %T", s))
	}
}

func (p *printer) print(text string) {
	p.js = append(p.js, text...)
}

func (p *printer) printNewline() {
	p.print("\n")
}
