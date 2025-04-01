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

	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}

func checkItem(
	env *Environment[tast.Type],
	typ tast.Type,
	item parser.IItemContext,
) (tast.Item, error) {
	line, col, text := extractPosData(item)

	switch i := item.(type) {
	case *parser.NoInitItemContext:
		varName := i.Ident().GetText()

		currentCtx, ok := env.Peek()
		if !ok {
			return nil, fmt.Errorf(
				"declaring variable outside of code blocks at %d:%d near '%s'",
				line, col, text,
			)
		}
		if currentCtx.Has(varName) {
			return nil, fmt.Errorf(
				"variable with name '%s' declared twice at %d:%d near '%s'",
				varName, line, col, text,
			)
		}

		(*currentCtx)[varName] = typ
		return tast.NewNoInitItem(varName, typ, line, col, text), nil

	case *parser.InitItemContext:
		varName := i.Ident().GetText()

		currentCtx, ok := env.Peek()
		if !ok {
			return nil, fmt.Errorf(
				"declaring variable outside of code blocks at %d:%d near '%s'",
				line, col, text,
			)
		}
		if currentCtx.Has(varName) {
			return nil, fmt.Errorf(
				"variable with name '%s' declared twice at %d:%d near '%s'",
				varName, line, col, text,
			)
		}

		(*currentCtx)[varName] = typ

		exp, err := inferExp(env, i.Exp())
		if err != nil {
			return nil, err
		}
		return tast.NewInitItem(varName, exp, typ, line, col, text), nil

	default:
		return nil, fmt.Errorf(
			"checkItem: unhandled item type %T at %d:%d near '%s'",
			i, line, col, text,
		)
	}
}
