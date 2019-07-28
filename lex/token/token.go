package token

import (
	"fmt"
	"strconv"
)

type Kind int
type Category int

const (
	Operator Category = iota
	None
)

const (
	Plus Kind = iota
	Minus
	Div
	Mul
	Mod
	Pipe
	LParen
	RParen
	Integer
	Identifier
	Assignment
	NewLine
	TernaryStart
	TernarySep
	EOF
	Invalid
)

func (k Kind) String() (s string) {
	switch k {
	case Plus:
		s = "+"
	case Div:
		s = "/"
	case Mul:
		s = "*"
	case Mod:
		s = "%"
	case Assignment:
		s = "="
	case Pipe:
		s = "|"
	case Identifier:
		s = "Id"
	case Integer:
		s = "int"
	case TernaryStart:
		s = "?"
	case TernarySep:
		s = ":"
	case EOF:
		s = "EOF"
	case NewLine:
		s = "\\n"
	default:
		s = strconv.Itoa(int(k))
	}
	return s
}

type Token struct {
	Kind    Kind
	Content string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s %s)", t.Kind, t.Content)
}
