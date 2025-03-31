package typechecker

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func checkStm(
	env *Environment[typedast.Type],
	stm parser.IStmContext,
) (typedast.Stm, error) {
	line, col, text := extractPosData(stm)
	switch s := stm.(type) {
	case *parser.ExpStmContext:
		inferredExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		return typedast.NewExpStm(inferredExp, line, col, text), nil

	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}
