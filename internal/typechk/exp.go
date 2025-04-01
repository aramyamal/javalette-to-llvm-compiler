package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func inferExp(
	env *Environment[tast.Type],
	exp parser.IExpContext,
) (tast.Exp, error) {
	line, col, text := extractPosData(exp)
	switch e := exp.(type) {
	case *parser.ParenExpContext:
		innerExp, err := inferExp(env, e.Exp())
		if err != nil {
			return nil, err
		}
		return tast.NewParenExp(innerExp, innerExp.Type(), line, col, text),
			nil
	default:
		return nil, fmt.Errorf(
			"inferExp: unhandled exp type %T at %d:%d near '%s'",
			e, line, col, text,
		)
	}
}

func promoteExp(exp tast.Exp, typ tast.Type) tast.Exp {
	if exp.Type() == tast.Int && typ == tast.Double {
		return tast.NewIntToDoubleExp(exp)
	}
	return exp
}
