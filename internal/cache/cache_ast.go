package cache

import (
	"github.com/demouth/esbundle/internal/js_ast"
	"github.com/demouth/esbundle/internal/js_parser"
	"github.com/demouth/esbundle/internal/logger"
)

type JSCache struct {
}

func (c *JSCache) Parse(source logger.Source) (js_ast.AST, bool) {

	ast := js_parser.Parse(source)

	// TODO
	return ast, true
}
