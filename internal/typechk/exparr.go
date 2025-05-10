package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) inferNewArrExp(
	e *parser.NewArrExpContext, line, col int, text string,
) (*tast.NewArrExp, error) {
	typ, err := toTastBaseType(e.BaseType())
	if err != nil {
		return nil, err
	}
	var indexExps []tast.Exp
	for i, idx := range e.AllArrayIndex() {
		typedExp, err := tc.inferExp(idx.Exp())
		if err != nil {
			return nil, err
		}
		if typedExp.Type() != tast.Int {
			return nil, fmt.Errorf(
				"array index at dimension %d must be of integer type at %d:%d "+
					" near %s", i, line, col, text,
			)
		}
		indexExps = append(indexExps, typedExp)
	}

	for range indexExps {
		typ = tast.Array(typ)
	}

	return tast.NewNewArrExp(indexExps, typ, line, col, text), nil
}

func (tc *TypeChecker) inferArrIndexExp(
	e *parser.ArrIndexExpContext, line, col int, text string,
) (*tast.ArrIndexExp, error) {
	exp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	typ, idxExps, err := tc.inferArrayIndexing(
		exp,
		e.AllArrayIndex(),
		line, col, text,
	)
	if err != nil {
		return nil, err
	}
	return tast.NewArrIndexExp(exp, idxExps, typ, line, col, text), nil
}

func (tc *TypeChecker) inferArrPostExp(
	e *parser.ArrPostExpContext, line, col int, text string,
) (*tast.ArrPostExp, error) {
	arrayExp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	typ, idxExps, err := tc.inferArrayIndexing(
		arrayExp,
		e.AllArrayIndex(),
		line, col, text,
	)
	if err != nil {
		return nil, err
	}

	if typ != tast.Int { //&& typ != tast.Double {
		return nil, fmt.Errorf(
			// "'++' or '--' operation can only be done on int or double at "+
			"'++' or '--' operation can only be done on int at "+
				"%d:%d near '%s'", line, col, text,
		)
	}
	op, err := extractIncDecOp(e.IncDecOp())
	if err != nil {
		return nil, fmt.Errorf("%w at %d:%d near %s", err, line, col, text)
	}
	return tast.NewArrPostExp(arrayExp, idxExps, op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferArrPreExp(
	e *parser.ArrPreExpContext, line, col int, text string,
) (*tast.ArrPreExp, error) {
	arrayExp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}

	typ, idxExps, err := tc.inferArrayIndexing(
		arrayExp,
		e.AllArrayIndex(),
		line, col, text,
	)
	if err != nil {
		return nil, err
	}

	if typ != tast.Int { //&& typ != tast.Double {
		return nil, fmt.Errorf(
			// "'++' or '--' operation can only be done on int or double at "+
			"'++' or '--' operation can only be done on int at "+
				"%d:%d near '%s'", line, col, text,
		)
	}
	op, err := extractIncDecOp(e.IncDecOp())
	if err != nil {
		return nil, fmt.Errorf("%w at %d:%d near %s", err, line, col, text)
	}
	return tast.NewArrPreExp(arrayExp, idxExps, op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferArrAssignExp(
	e *parser.ArrAssignExpContext, line, col int, text string,
) (*tast.ArrAssignExp, error) {

	arrExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	typ, idxExps, err := tc.inferArrayIndexing(
		arrExp,
		e.AllArrayIndex(),
		line, col, text,
	)
	if err != nil {
		return nil, err
	}
	assExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}
	if !isConvertible(typ, assExp.Type()) {
		return nil, fmt.Errorf(
			"array access assignment with wrong type at %d:%d near %s, "+
				"expected type %s but got type %s",
			line, col, text, typ.String(), assExp.Type().String(),
		)
	}
	return tast.NewArrAssignExp(arrExp, idxExps, assExp, typ, line, col, text),
		nil
}

func (tc *TypeChecker) inferArrayIndexing(
	arrayExp tast.Exp,
	allArrayIndexContext []parser.IArrayIndexContext,
	line, col int, text string,
) (tast.Type, []tast.Exp, error) {

	currentType := arrayExp.Type()
	var idxExps []tast.Exp

	for i, idx := range allArrayIndexContext {
		idxExp, err := tc.inferExp(idx.Exp())
		if err != nil {
			return nil, nil, err
		}

		if idxExp.Type() != tast.Int {
			return nil, nil, fmt.Errorf(
				"array index access at dimension %d must be integer type at "+
					"%d:%d near %s", i, line, col, text,
			)
		}
		idxExps = append(idxExps, idxExp)

		arrType, ok := currentType.(*tast.ArrayType)
		if !ok {
			return nil, nil, fmt.Errorf(
				"array index mismatch at dimension %d at %d:%d near %s "+
					", expected array type but got %s instead",
				i, line, col, text, currentType.String(),
			)
		}
		currentType = arrType.Elem
	}
	return currentType, idxExps, nil
}
