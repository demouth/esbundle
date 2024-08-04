package graph

import (
	"github.com/demouth/esbundle/internal/ast"
	"github.com/demouth/esbundle/internal/js_ast"
)

type InputFile struct {
	Repr InputFileRepr
}

type InputFileRepr interface {
	ImportRecords() *[]ast.ImportRecord
}

type JSRepr struct {
	AST js_ast.AST
}

func (f *JSRepr) ImportRecords() *[]ast.ImportRecord {
	// TODO
	return nil
}

type OutputFile struct {
	AbsPath  string
	Contents []byte
}
