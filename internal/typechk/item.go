package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) checkItem(
	typ tast.Type,
	item parser.IItemContext,
) (tast.Item, error) {
	line, col, text := extractPosData(item)

	switch i := item.(type) {
	case *parser.NoInitItemContext:
		return tc.checkNoInitItem(typ, i, line, col, text)
	case *parser.InitItemContext:
		return tc.checkInitItem(typ, i, line, col, text)
	default:
		return nil, fmt.Errorf(
			"checkItem: unhandled item type %T at %d:%d near '%s'",
			i, line, col, text,
		)
	}
}

func (tc *TypeChecker) checkNoInitItem(
	typ tast.Type,
	i *parser.NoInitItemContext,
	line, col int, text string,
) (*tast.NoInitItem, error) {
	varName := i.Ident().GetText()
	currentCtx, ok := tc.env.Peek()
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
}

func (tc *TypeChecker) checkInitItem(
	typ tast.Type,
	i *parser.InitItemContext,
	line, col int, text string,
) (*tast.InitItem, error) {
	varName := i.Ident().GetText()

	currentCtx, ok := tc.env.Peek()
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

	typedExp, err := tc.inferExp(i.Exp())
	if err != nil {
		return nil, err
	}

	if !isConvertible(typ, typedExp.Type()) {
		return nil, fmt.Errorf(
			"cannot convert from %s to %s at %d:%d near '%s'",
			typedExp.Type().String(), typ.String(), line, col, text,
		)
	}

	return tast.NewInitItem(varName, typedExp, typ, line, col, text), nil
}
