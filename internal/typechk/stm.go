package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) checkStm(stm parser.IStmContext) (tast.Stm, error) {
	line, col, text := extractPosData(stm)
	switch s := stm.(type) {
	case *parser.ExpStmContext:
		return tc.checkExpStm(s, line, col, text)
	case *parser.DeclsStmContext:
		return tc.checkDeclsStm(s, line, col, text)
	case *parser.ReturnStmContext:
		return tc.checkReturnStm(s, line, col, text)
	case *parser.VoidReturnStmContext:
		return tc.checkVoidReturnStm(line, col, text)
	case *parser.ForEachStmContext:
		return tc.checkForEachStm(s, line, col, text)
	case *parser.WhileStmContext:
		return tc.checkWhileStm(s, line, col, text)
	case *parser.BlockStmContext:
		return tc.checkBlockStm(s, line, col, text)
	case *parser.IfStmContext:
		return tc.checkIfStm(s, line, col, text)
	case *parser.BlankStmContext:
		return tast.NewBlankStm(line, col, text), nil
	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}
