package parser

import (
	"fmt"
	"strconv"

	"github.com/MatiasLyyra/JabaScript/lex/token"
)

type parser struct {
	tokens []token.Token
	idx    int
}

func (p *parser) peek() token.Token {
	if p.idx >= len(p.tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.tokens[p.idx]
}
func (p *parser) peekAt(offet int) token.Token {
	if offet+p.idx < len(p.tokens) {
		return p.tokens[offet+p.idx]
	}
	return token.Token{Kind: token.EOF}
}
func (p *parser) consume() token.Token {
	prevToken := p.peek()
	p.idx++
	return prevToken
}

func (p *parser) accept(kind token.Kind) bool {
	return p.peek().Kind == kind
}

func (p *parser) require(kind token.Kind) error {
	if p.peek().Kind == kind {
		p.consume()
		return nil
	}
	return fmt.Errorf("expected token %s found %s", kind, p.peek().Kind)
}

func (p *parser) isWhiteSpace() bool {
	return p.accept(token.EOF) || p.accept(token.NewLine)
}

func (p *parser) assignmentExpression() (Expression, error) {
	if p.accept(token.Identifier) && p.peekAt(1).Kind == token.Assignment {
		exp := AssignmentExpression{id: p.consume().Content}
		// Consume the = operator
		p.consume()
		if p.isWhiteSpace() {
			return nil, fmt.Errorf("unexpected \\n or EOF, expected expression")
		}
		val, err := p.assignmentExpression()
		exp.val = val
		return exp, err
	}
	return p.addExpression()
}

func (p *parser) addExpression() (Expression, error) {
	tree, err := p.mulExpression()
	if err != nil {
		return nil, err
	}
	for p.accept(token.Plus) || p.accept(token.Minus) {
		op := p.consume()
		addExp := BinaryExpression{lExp: tree, op: op.Content}
		if p.isWhiteSpace() {
			return nil, fmt.Errorf("unexpected \\n or EOF, expected expression")
		}
		rExp, err := p.mulExpression()
		if err != nil {
			return nil, err
		}
		addExp.rExp = rExp
		tree = addExp
	}
	return tree, nil
}
func (p *parser) mulExpression() (Expression, error) {
	tree, err := p.function()
	if err != nil {
		return nil, err
	}
	for p.accept(token.Mul) ||
		p.accept(token.Div) ||
		p.accept(token.Mod) {
		op := p.consume().Content
		exp := BinaryExpression{op: op, lExp: tree}
		if p.isWhiteSpace() {
			return nil, fmt.Errorf("unexpected \\n or EOF, expected expression")
		}
		rExp, err := p.function()
		if err != nil {
			return nil, err
		}
		exp.rExp = rExp
		tree = exp
	}
	return tree, nil
}

func (p *parser) unary() (Expression, error) {
	uExp := UnaryExpression{sign: 1}
	if p.accept(token.Minus) {
		p.consume()
		uExp.sign = -1
	}
	if p.accept(token.Integer) {
		val, _ := strconv.Atoi(p.consume().Content)
		uExp.val = IntegerExpression(val)
	} else if p.accept(token.Identifier) {
		idName := p.consume().Content
		uExp.val = IdentifierExpression(idName)
	} else if p.accept(token.LParen) {
		p.consume()
		val, err := p.assignmentExpression()
		if err != nil {
			return nil, err
		}
		uExp.val = val
		err = p.require(token.RParen)
		if err != nil {
			return nil, err
		}
	} else if p.accept(token.Pipe) {
		p.consume()
		fnExp := FnDefExpression{
			arguments: make([]string, 0),
		}
		for p.accept(token.Identifier) {
			fnExp.arguments = append(fnExp.arguments, p.consume().Content)
		}
		err := p.require(token.Pipe)
		if err != nil {
			return nil, err
		}
		body, err := p.assignmentExpression()
		if err != nil {
			return nil, err
		}
		fnExp.body = body
		uExp.val = fnExp
	} else {
		return nil, fmt.Errorf("unexpected token %s", p.consume().Kind)
	}
	return uExp, nil
}
func (p *parser) function() (Expression, error) {
	fn, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.accept(token.LParen) {
		fnCallExp := FnCallExpression{
			fn:        fn,
			arguments: make([]Expression, 0),
		}
		p.consume()
		for !p.accept(token.RParen) {
			if p.isWhiteSpace() {
				break
			}
			arg, err := p.addExpression()
			if err != nil {
				return nil, err
			}
			fnCallExp.arguments = append(fnCallExp.arguments, arg)
			fn = fnCallExp
		}
		err := p.require(token.RParen)
		if err != nil {
			return nil, err
		}
	}

	return fn, nil
}
func (p *parser) program() (Expression, error) {
	val, err := p.assignmentExpression()
	if err != nil {
		return nil, err
	}
	err = p.require(token.NewLine)
	return val, err
}
func Parse(tokens []token.Token) (Expression, error) {
	p := parser{tokens: tokens}
	return p.program()
}
