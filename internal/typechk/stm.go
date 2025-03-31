package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func checkStm(
	env *Environment[tast.Type],
	stm parser.IStmContext,
) (tast.Stm, error) {
	line, col, text := extractPosData(stm)
	switch s := stm.(type) {
	case *parser.ExpStmContext:
		inferredExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		return tast.NewExpStm(inferredExp, line, col, text), nil

	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}
