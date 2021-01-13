package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type MinusAssign struct {
	l Node
	r Node
}

func NewMinusAssign(l, r Node) Node {
	return MinusAssign{l, r}
}

func (m MinusAssign) Eval(env *obj.Env) obj.Object {
	var name string
	var left = m.l.Eval(env)
	var right = m.r.Eval(env)

	if ident, ok := m.l.(Identifier); ok {
		name = ident.String()
	} else {
		return obj.NewError("cannot assign to literal")
	}

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if !assertTypes(left, obj.INT, obj.FLOAT) {
		return obj.NewError("unsupported operator '-=' for type %v", left.Type())
	}
	if !assertTypes(right, obj.INT, obj.FLOAT) {
		return obj.NewError("unsupported operator '-=' for type %v", right.Type())
	}

	switch {
	case assertTypes(left, obj.INT) && assertTypes(right, obj.INT):
		l := left.(*obj.Integer).Val()
		r := right.(*obj.Integer).Val()
		env.Set(name, obj.NewInteger(l-r))

	case assertTypes(left, obj.FLOAT, obj.INT) && assertTypes(right, obj.FLOAT, obj.INT):
		left, right = toFloat(left, right)
		l := left.(*obj.Float).Val()
		r := right.(*obj.Float).Val()
		env.Set(name, obj.NewFloat(l-r))

	default:
		return obj.NewError(
			"invalid operation %v -= %v (wrong types %v and %v)",
			left, right, left.Type(), right.Type(),
		)
	}

	return obj.NullObj
}

func (m MinusAssign) String() string {
	return fmt.Sprintf("(%v -= %v)", m.l, m.r)
}
