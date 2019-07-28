package lex

import (
	"bytes"
	"fmt"
	"io"
	"text/scanner"
	"unicode"

	"github.com/MatiasLyyra/JabaScript/lex/token"
)

type lexer struct {
	scanner scanner.Scanner
	buf     *bytes.Buffer
	current rune
}

func (lex *lexer) skipWhiteSpace() {
	for unicode.IsSpace(lex.current) && lex.current != '\n' {
		lex.advance(false)
	}
}

func (lex *lexer) advance(save bool) {
	if lex.current != scanner.EOF {
		char := lex.scanner.Next()
		if save {
			lex.buf.WriteRune(lex.current)
		}
		lex.current = char
	}
}

func (lex *lexer) isDigit() bool {
	return unicode.IsDigit(lex.current)
}

func (lex *lexer) isLetter() bool {
	return unicode.IsLetter(lex.current)
}

func (lex *lexer) isOperator() bool {
	return lex.current == '/' ||
		lex.current == '*' ||
		lex.current == '%' ||
		lex.current == '-' ||
		lex.current == '+' ||
		lex.current == '=' ||
		lex.current == '(' ||
		lex.current == ')' ||
		lex.current == '?' ||
		lex.current == ':' ||
		lex.current == '|'
}
func (lex *lexer) scanDigit() token.Token {
	for lex.isDigit() {
		lex.advance(true)
	}

	return lex.makeToken(token.Integer)
}

func (lex *lexer) scanIdentifier() token.Token {
	for lex.isLetter() || lex.isDigit() || lex.current == '_' {
		lex.advance(true)
	}

	return lex.makeToken(token.Identifier)
}

func (lex *lexer) makeToken(kind token.Kind) token.Token {
	token := token.Token{Content: lex.buf.String(), Kind: kind}
	lex.buf.Reset()
	return token
}

func (lex *lexer) scanOperator() token.Token {
	var kind token.Kind
	switch lex.current {
	case '*':
		kind = token.Mul
	case '+':
		kind = token.Plus
	case '-':
		kind = token.Minus
	case '(':
		kind = token.LParen
	case ')':
		kind = token.RParen
	case '|':
		kind = token.Pipe
	case '/':
		kind = token.Div
	case '?':
		kind = token.TernaryStart
	case ':':
		kind = token.TernarySep
	case '%':
		kind = token.Mod
	case '=':
		kind = token.Assignment
	}
	lex.advance(true)
	return lex.makeToken(kind)
}

func (lex *lexer) isNewLine() bool {
	return lex.current == '\n'
}

func (lex *lexer) scanNewLine() token.Token {
	for lex.isNewLine() {
		lex.advance(false)
	}
	return lex.makeToken(token.NewLine)
}
func (lex *lexer) next() (token.Token, error) {
	lex.skipWhiteSpace()
	if lex.current == scanner.EOF {
		lex.advance(true)
		return lex.makeToken(token.EOF), nil
	}
	if lex.isNewLine() {
		return lex.scanNewLine(), nil
	}
	if lex.isDigit() {
		return lex.scanDigit(), nil
	}
	if lex.isLetter() {
		return lex.scanIdentifier(), nil
	}
	if lex.isOperator() {
		return lex.scanOperator(), nil
	}
	return token.Token{}, fmt.Errorf("invalid token %s", lex.buf.String())
}

func Tokenize(r io.Reader) ([]token.Token, error) {
	tokens := make([]token.Token, 0)
	lex := &lexer{
		buf: bytes.NewBufferString(""),
	}
	lex.scanner.Init(r)
	var (
		err error
		t   token.Token
	)
	lex.advance(false)
	for lex.current != scanner.EOF {
		t, err = lex.next()
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, t)
	}
	return tokens, err
}
