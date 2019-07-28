package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Type int

const (
	Integer Type = iota
	Function
	Closure
)
const maxStackSize = 500

type Value interface{}
type ContextValue struct {
	Val  Value
	Type Type
}

type ClosureContext struct {
	vars map[string]ContextValue
	fn   *FnDefExpression
}

func (c ContextValue) String() string {
	switch c.Type {
	case Integer:
		return strconv.Itoa(c.Val.(int))
	case Function:
		return "[Function]"
	case Closure:
		return "[Function]"
	default:
		return fmt.Sprintf("Unknown val: (%d %s)", c.Type, c.Val)
	}
}

type Context struct {
	Vars       map[string]ContextValue
	stackCount int
}

func (c *Context) merge(vars map[string]ContextValue) {
	for key, val := range vars {
		if _, ok := c.Vars[key]; !ok {
			c.Vars[key] = val
		}
	}
}

func NewContext() *Context {
	ctx := &Context{
		Vars: make(map[string]ContextValue),
	}
	return ctx
}

type Expression interface {
	Eval(*Context) (ContextValue, error)
}

type FnDefExpression struct {
	arguments []string
	body      Expression
}

func (exp FnDefExpression) String() string {
	return fmt.Sprintf("(fn |%s| (%s))", strings.Join(exp.arguments, ","), exp.body)
}

func (exp FnDefExpression) Eval(ctx *Context) (ContextValue, error) {
	return ContextValue{Type: Function, Val: exp}, nil
}

type FnCallExpression struct {
	arguments []Expression
	fn        Expression
}

func (exp FnCallExpression) String() string {
	args := make([]string, 0, len(exp.arguments))
	for _, arg := range exp.arguments {
		args = append(args, fmt.Sprint(arg))
	}
	return fmt.Sprintf("(fn call %s (%s))", strings.Join(args, " "), exp.fn)
}

func (exp FnCallExpression) Eval(ctx *Context) (ContextValue, error) {
	fnCtx, err := exp.fn.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	if fnCtx.Type != Function && fnCtx.Type != Closure {
		return ContextValue{}, fmt.Errorf("expression does not evaluate to a function")
	}
	vars := make(map[string]ContextValue)
	var fn *FnDefExpression
	if fnCtx.Type == Function {
		fnVal := fnCtx.Val.(FnDefExpression)
		fn = &fnVal
	} else {
		closure := fnCtx.Val.(ClosureContext)
		fn = closure.fn
		vars = closure.vars
	}

	if len(fn.arguments) != len(exp.arguments) {
		return ContextValue{}, fmt.Errorf("incorrent amount of arguments for \"%s\", expected %d got %d", exp.fn, len(fn.arguments), len(exp.arguments))
	}
	ctx.stackCount++
	if ctx.stackCount > maxStackSize {
		ctx.stackCount = 0
		return ContextValue{}, fmt.Errorf("max stack size %d exceeded", maxStackSize)
	}
	for i := 0; i < len(exp.arguments); i++ {
		argVal, err := exp.arguments[i].Eval(ctx)
		if err != nil {
			return ContextValue{}, err
		}
		vars[fn.arguments[i]] = argVal
	}
	oldVars := ctx.Vars
	ctx.Vars = vars
	ctx.merge(oldVars)
	fnVal, err := fn.body.Eval(ctx)
	if fnVal.Type == Function {
		closureFn := fnVal.Val.(FnDefExpression)
		fnVal = ContextValue{
			Val: ClosureContext{
				vars: ctx.Vars,
				fn:   &closureFn,
			},
			Type: Closure,
		}
	}
	ctx.Vars = oldVars
	ctx.stackCount--
	return fnVal, err
}

type AssignmentExpression struct {
	id  string
	val Expression
}

func (exp AssignmentExpression) String() string {
	return fmt.Sprintf("(= %s %s)", exp.id, exp.val)
}

func (exp AssignmentExpression) Eval(ctx *Context) (ContextValue, error) {
	val, err := exp.val.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	ctx.Vars[exp.id] = val
	return val, nil
}

type TernaryExpression struct {
	tExp Expression
	fExp Expression
	cond Expression
}

func (exp TernaryExpression) String() string {
	return fmt.Sprintf("(? %s %s %s)", exp.cond, exp.tExp, exp.fExp)
}

func (exp TernaryExpression) Eval(ctx *Context) (ContextValue, error) {
	valCtx, err := exp.cond.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	if valCtx.Type != Integer {
		return ContextValue{}, fmt.Errorf("cannot evaluate condition on non integer value")
	}
	val := valCtx.Val.(int)

	if val != 0 {
		return exp.tExp.Eval(ctx)
	}
	return exp.fExp.Eval(ctx)
}

type BinaryExpression struct {
	lExp Expression
	rExp Expression
	op   string
}

func (exp BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", exp.op, exp.lExp, exp.rExp)
}

func (exp BinaryExpression) Eval(ctx *Context) (ContextValue, error) {
	lCtxVal, err := exp.lExp.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	if lCtxVal.Type != Integer {
		return ContextValue{}, fmt.Errorf("%s cannot be applied to expression", exp.op)
	}
	rCtxVal, err := exp.rExp.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	if rCtxVal.Type != Integer {
		return ContextValue{}, fmt.Errorf("%s cannot be applied to expression", exp.op)
	}
	var lVal, rVal int = lCtxVal.Val.(int), rCtxVal.Val.(int)

	var val int
	switch exp.op {
	case "%":
		val = lVal % rVal
	case "*":
		val = lVal * rVal
	case "/":
		if rVal == 0 {
			return ContextValue{}, fmt.Errorf("division by zero")
		}
		val = lVal / rVal
	case "+":
		val = lVal + rVal
	case "-":
		val = lVal - rVal
	}
	return ContextValue{Val: val, Type: Integer}, nil
}

type UnaryExpression struct {
	sign int
	val  Expression
}

func (exp UnaryExpression) Eval(ctx *Context) (ContextValue, error) {
	ctxVal, err := exp.val.Eval(ctx)
	if err != nil {
		return ContextValue{}, err
	}
	if ctxVal.Type != Integer && exp.sign != 1 {
		return ContextValue{}, fmt.Errorf("cannot apply - to non integer")
	}
	if exp.sign == -1 {
		return ContextValue{Val: -ctxVal.Val.(int), Type: Integer}, nil
	}
	return ctxVal, nil
}

func (exp UnaryExpression) String() string {
	if exp.sign == -1 {
		return fmt.Sprintf("-%s", exp.val)
	}
	return fmt.Sprintf("%s", exp.val)
}

type IntegerExpression int

func (exp IntegerExpression) Eval(ctx *Context) (ContextValue, error) {
	return ContextValue{Val: int(exp), Type: Integer}, nil
}

func (i IntegerExpression) String() string {
	return strconv.Itoa(int(i))
}

type IdentifierExpression string

func (exp IdentifierExpression) Eval(ctx *Context) (ContextValue, error) {
	if val, ok := ctx.Vars[string(exp)]; ok {
		return val, nil
	}
	return ContextValue{}, fmt.Errorf("variable \"%s\" not defined", exp)
}

func (id IdentifierExpression) String() string {
	return string(id)
}
