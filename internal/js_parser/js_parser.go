package js_parser

import (
	"github.com/demouth/esbundle/internal/js_ast"
	"github.com/demouth/esbundle/internal/js_lexer"
	"github.com/demouth/esbundle/internal/logger"
)

type parser struct {
	source logger.Source
	lexer  js_lexer.Lexer
}

func Parse(source logger.Source) (result js_ast.AST) {
	p := newParser(source, js_lexer.NewLexer(source))
	stmts := p.parseStmtsUpTo(js_lexer.TEndOfFile)

	var parts []js_ast.Part

	parts = p.appendPart(parts, stmts)

	result = p.toAST(parts)
	return
}

func newParser(source logger.Source, lexer js_lexer.Lexer) *parser {
	return &parser{
		source: source,
		lexer:  lexer,
	}
}

func (p *parser) parseStmtsUpTo(end js_lexer.T) []js_ast.Stmt {
	stmts := []js_ast.Stmt{}
	for {
		if p.lexer.Token == end {
			break
		}
		stmt := p.parseStmt()
		stmts = append(stmts, stmt)
	}
	return stmts
}

func (p *parser) parseStmt() js_ast.Stmt {
	switch p.lexer.Token {
	case js_lexer.TSemicolon:
		p.lexer.Next()
		return js_ast.Stmt{
			Data: js_ast.SEmptyShared,
		}
	}
	return js_ast.Stmt{}
}

func (p *parser) toAST(parts []js_ast.Part) js_ast.AST {
	return js_ast.AST{
		Parts: parts,
	}
}

func (p *parser) appendPart(parts []js_ast.Part, stmts []js_ast.Stmt) []js_ast.Part {
	part := js_ast.Part{
		Stmts: p.visitStmtsAndPrependTempRefs(stmts),
	}
	parts = append(parts, part)
	return parts
}

func (p *parser) visitStmtsAndPrependTempRefs(stmts []js_ast.Stmt) []js_ast.Stmt {
	return stmts
}
