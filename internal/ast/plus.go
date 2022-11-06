package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type Plus struct {
	l   Node
	r   Node
	pos int
}

func NewPlus(l, r Node, pos int) Node {
	return Plus{
		l:   l,
		r:   r,
		pos: pos,
	}
}

func (p Plus) Eval(env *obj.Env) obj.Object {
	var (
		left  = obj.Unwrap(p.l.Eval(env))
		right = obj.Unwrap(p.r.Eval(env))
	)

	if takesPrecedence(left) {
		return left
	}
	if takesPrecedence(right) {
		return right
	}

	if !obj.AssertTypes(left, obj.IntType, obj.FloatType, obj.StringType) {
		return obj.NewError("unsupported operator '+' for type %v", left.Type())
	}
	if !obj.AssertTypes(right, obj.IntType, obj.FloatType, obj.StringType) {
		return obj.NewError("unsupported operator '+' for type %v", right.Type())
	}

	switch {
	case obj.AssertTypes(left, obj.StringType) && obj.AssertTypes(right, obj.StringType):
		l := left.(obj.String)
		r := right.(obj.String)
		return obj.String(l + r)

	case obj.AssertTypes(left, obj.IntType) && obj.AssertTypes(right, obj.IntType):
		l := left.(obj.Integer)
		r := right.(obj.Integer)
		return obj.Integer(l + r)

	case obj.AssertTypes(left, obj.FloatType, obj.IntType) && obj.AssertTypes(right, obj.FloatType, obj.IntType):
		left, right = obj.ToFloat(left, right)
		l := left.(obj.Float)
		r := right.(obj.Float)
		return obj.Float(l + r)

	default:
		return obj.NewError(
			"invalid operation %v + %v (wrong types %v and %v)",
			left, right, left.Type(), right.Type(),
		)
	}
}

func (p Plus) String() string {
	return fmt.Sprintf("(%v + %v)", p.l, p.r)
}

func (p Plus) Compile(c *compiler.Compiler) (position int, err error) {
	if p.IsConstExpression() {
		o := p.Eval(nil)
		if e, ok := o.(*obj.Error); ok {
			return 0, compiler.NewError(p.pos, string(*e))
		}
		position = c.Emit(code.OpConstant, c.AddConstant(o))
		c.Bookmark(p.pos)
		return
	}

	if position, err = p.l.Compile(c); err != nil {
		return
	}
	if position, err = p.r.Compile(c); err != nil {
		return
	}
	position = c.Emit(code.OpAdd)
	c.Bookmark(p.pos)
	return
}

func (p Plus) IsConstExpression() bool {
	return p.l.IsConstExpression() && p.r.IsConstExpression()
}
