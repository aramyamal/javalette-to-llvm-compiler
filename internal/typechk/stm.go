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

	case *parser.DeclsStmContext:
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

	case *parser.ReturnStmContext:
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
			"illegal conversion in return. Expected %s, "+
				"but got %s instead",
			returnType.String(), expType.String(),
		)
	case *parser.VoidReturnStmContext:
		returnType := tc.env.ReturnType()
		if isConvertible(returnType, types.Void) {
			return tast.NewVoidReturnStm(line, col, text), nil
		}
		return nil, fmt.Errorf(
			"illegal conversion in return. Expected %s, "+
				"but got %s instead",
			returnType.String(), types.Void.String(),
		)
	case *parser.WhileStmContext:
		typedExp, err := tc.inferExp(s.Exp())
		if err != nil {
			return nil, err
		}
		if typedExp.Type() != types.Bool {
			return nil, fmt.Errorf(
				"expression in while-loop does not have type bool at %d:%d "+
					"near '%s'",
				line, col, text,
			)
		}
		tc.env.EnterContext()
		typedStm, err := tc.checkStm(s.Stm())
		if err != nil {
			return nil, err
		}
		tc.env.ExitContext()
		return tast.NewWhileStm(typedExp, typedStm, line, col, text), nil

	case *parser.BlockStmContext:
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

	case *parser.IfStmContext:
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

	case *parser.BlankStmContext:
		return tast.NewBlankStm(line, col, text), nil

	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}
