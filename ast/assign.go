package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type Assign struct {
	l Node
	r Node
}

func NewAssign(l, r Node) Node {
	return Assign{l, r}
}

func (a Assign) Eval(env *obj.Env) obj.Object {
	if left, ok := a.l.(Identifier); ok {
		right := obj.Unwrap(a.r.Eval(env))
		if takesPrecedence(right) {
			return right
		}
		return env.Set(left.String(), right)
	}

	left := a.l.Eval(env)
	if s, ok := left.(obj.Setter); ok {
		right := obj.Unwrap(a.r.Eval(env))
		if takesPrecedence(right) {
			return right
		}
		return s.Set(right)
	}

	return obj.NewError("cannot assign to literal")
}

func (a Assign) String() string {
	return fmt.Sprintf("(%v = %v)", a.l, a.r)
}
