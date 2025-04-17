package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
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
		return tc.checkVoidReturnStm(s, line, col, text)
	case *parser.WhileStmContext:
		return tc.checkWhileStm(s, line, col, text)
	case *parser.BlockStmContext:
		return tc.checkBlockStm(s, line, col, text)
	case *parser.IfStmContext:
		return tc.checkIfStm(s, line, col, text)
	case *parser.BlankStmContext:
		return tc.checkBlankStm(s, line, col, text)
	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}

func (tc *TypeChecker) checkExpStm(
	s *parser.ExpStmContext, line, col int, text string,
) (*tast.ExpStm, error) {
	typedExp, err := tc.inferExp(s.Exp())
	if err != nil {
		return nil, err
	}

	if !typedExp.HasSideEffect() {
		return nil, fmt.Errorf(
			"expression statement has no effect at %d:%d near '%s'",
			line, col, text,
		)
	}
	return tast.NewExpStm(typedExp, line, col, text), nil
}

func (tc *TypeChecker) checkDeclsStm(
	s *parser.DeclsStmContext, line, col int, text string,
) (*tast.DeclsStm, error) {
	typ, err := toIrType(s.Type_())
	if err != nil {
		return nil, err
	}
	if typ == types.Void {
		return nil, fmt.Errorf(
			"variable declaration of type void at %d:%d near '%s'",
			line, col, text,
		)
	}
	items := []tast.Item{}
	for _, item := range s.AllItem() {
		typedItem, err := tc.checkItem(typ, item)
		if err != nil {
			return nil, err
		}
		items = append(items, typedItem)
	}
	return tast.NewDeclsStm(items, line, col, text), nil
}

func (tc *TypeChecker) checkReturnStm(
	s *parser.ReturnStmContext, line, col int, text string,
) (*tast.ReturnStm, error) {
	typedExp, err := tc.inferExp(s.Exp())
	if err != nil {
		return nil, err
	}
	returnType := tc.env.ReturnType()
	expType := typedExp.Type()
	if isConvertible(returnType, expType) {
		return tast.NewReturnStm(
			returnType,
			promoteExp(typedExp, returnType),
			line, col, text,
		), nil
	}
	return nil, fmt.Errorf(
		"illegal conversion in return. Expected %s, but got %s instead",
		returnType.String(), expType.String(),
	)
}

func (tc *TypeChecker) checkVoidReturnStm(
	s *parser.VoidReturnStmContext, line, col int, text string,
) (*tast.VoidReturnStm, error) {
	returnType := tc.env.ReturnType()
	if isConvertible(returnType, types.Void) {
		return tast.NewVoidReturnStm(line, col, text), nil
	}
	return nil, fmt.Errorf(
		"illegal conversion in return. Expected %s, but got %s instead",
		returnType.String(), types.Void.String(),
	)
}

func (tc *TypeChecker) checkWhileStm(
	s *parser.WhileStmContext, line, col int, text string,
) (*tast.WhileStm, error) {
	typedExp, err := tc.inferExp(s.Exp())
	if err != nil {
		return nil, err
	}
	if typedExp.Type() != types.Bool {
		return nil, fmt.Errorf(
			"expression in while-loop does not have type bool at %d:%d "+
				"near '%s'", line, col, text,
		)
	}
	tc.env.EnterContext()
	typedStm, err := tc.checkStm(s.Stm())
	if err != nil {
		return nil, err
	}
	tc.env.ExitContext()
	return tast.NewWhileStm(typedExp, typedStm, line, col, text), nil
}

func (tc *TypeChecker) checkBlockStm(
	s *parser.BlockStmContext, line, col int, text string,
) (*tast.BlockStm, error) {
	tc.env.EnterContext()
	typedStms := []tast.Stm{}
	for _, stm := range s.AllStm() {
		typedStm, err := tc.checkStm(stm)
		if err != nil {
			return nil, err
		}
		typedStms = append(typedStms, typedStm)
	}
	tc.env.ExitContext()
	return tast.NewBlockStm(typedStms, line, col, text), nil
}

func (tc *TypeChecker) checkIfStm(
	s *parser.IfStmContext, line, col int, text string,
) (*tast.IfStm, error) {
	typedExp, err := tc.inferExp(s.Exp())
	if err != nil {
		return nil, err
	}
	if typedExp.Type() != types.Bool {
		return nil, fmt.Errorf(
			"if else expression does not have type bool at %d:%d near '%s'",
			line, col, text,
		)
	}

	tc.env.EnterContext()
	thenStm, err := tc.checkStm(s.Stm(0))
	if err != nil {
		return nil, err
	}
	tc.env.ExitContext()

	var elseStm tast.Stm = nil
	if len(s.AllStm()) > 1 {
		tc.env.EnterContext()
		elseStm, err = tc.checkStm(s.Stm(1))
		if err != nil {
			return nil, err
		}
		tc.env.ExitContext()
	}
	return tast.NewIfStm(typedExp, thenStm, elseStm, line, col, text), nil
}

func (tc *TypeChecker) checkBlankStm(
	s *parser.BlankStmContext, line, col int, text string,
) (*tast.BlankStm, error) {
	return tast.NewBlankStm(line, col, text), nil
}
