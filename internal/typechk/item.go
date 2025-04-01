package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

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
