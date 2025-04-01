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

	case *parser.DeclsStmContext:
		typ, err := toAstType(s.Type_())
		if err != nil {
			return nil, err
		}
		if typ == tast.Void {
			return nil, fmt.Errorf(
				"variable declaration of type void at %d:%d near '%s'",
				line, col, text,
			)
		}
		items := []tast.Item{}
		for _, item := range s.AllItem() {
			typedItem, err := checkItem(env, typ, item)
			if err != nil {
				return nil, err
			}
			items = append(items, typedItem)
		}
		return tast.NewDeclsStm(items, line, col, text), nil

	case *parser.ReturnStmContext:
		typedExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		returnType := env.ReturnType()
		expType := typedExp.Type()
		if isConvertible(returnType, expType) {
			return tast.NewReturnStm(
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
		returnType := env.ReturnType()
		if isConvertible(returnType, tast.Void) {
			return tast.NewVoidReturnStm(line, col, text), nil
		}
		return nil, fmt.Errorf(
			"illegal conversion in return. Expected %s, "+
				"but got %s instead",
			returnType.String(), tast.Void.String(),
		)
	case *parser.WhileStmContext:
		typedExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		if typedExp.Type() != tast.Void {
			return nil, fmt.Errorf(
				"expression in while-loop does not have type bool at %d:%d "+
					"near '%s'",
				line, col, text,
			)
		}
		env.EnterContext()
		typedStm, err := checkStm(env, s.Stm())
		if err != nil {
			return nil, err
		}
		env.ExitContext()
		return tast.NewWhileStm(typedExp, typedStm, line, col, text), nil

	case *parser.BlockStmContext:
		env.EnterContext()
		typedStms := []tast.Stm{}
		for _, stm := range s.AllStm() {
			typedStm, err := checkStm(env, stm)
			if err != nil {
				return nil, err
			}
			typedStms = append(typedStms, typedStm)
		}
		env.ExitContext()
		return tast.NewBlockStm(typedStms, line, col, text), nil

	case *parser.IfStmContext:
		typedExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		if typedExp.Type() != tast.Bool {
			return nil, fmt.Errorf(
				"if else expression does not have type bool at %d:%d near '%s'",
				line, col, text,
			)
		}

		env.EnterContext()
		thenStm, err := checkStm(env, s.Stm(0))
		if err != nil {
			return nil, err
		}
		env.ExitContext()

		var elseStm tast.Stm = nil
		if len(s.AllStm()) > 1 {
			env.EnterContext()
			elseStm, err = checkStm(env, s.Stm(1))
			if err != nil {
				return nil, err
			}
			env.ExitContext()
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
