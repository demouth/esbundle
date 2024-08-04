package js_lexer

import (
	"unicode/utf8"

	"github.com/demouth/esbundle/internal/logger"
)

type T uint8

const (
	TEndOfFile T = iota
	TSemicolon
)

type Lexer struct {
	source logger.Source

	codePoint rune
	current   int
	start     int
	end       int
	Token     T
}

func NewLexer(source logger.Source) Lexer {
	l := Lexer{
		source: source,
	}
	l.step()
	l.Next()
	return l
}

func (lexer *Lexer) step() {
	codePoint, width := utf8.DecodeRuneInString(lexer.source.Contents[lexer.current:])
	if width == 0 {
		codePoint = -1
	}
	lexer.codePoint = codePoint
	lexer.end = lexer.current
	lexer.current += width
}

func (lexer *Lexer) Next() {
	for {
		lexer.start = lexer.end
		lexer.Token = 0
		switch lexer.codePoint {
		case -1:
			lexer.Token = TEndOfFile
		case ';':
			lexer.step()
			lexer.Token = TSemicolon
		}
		return
	}
}
