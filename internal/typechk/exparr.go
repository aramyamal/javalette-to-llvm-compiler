package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) inferNewArrExp(
	e *parser.NewArrExpContext, line, col int, text string,
) (*tast.NewArrExp, error) {
	typ, err := tc.toTastBaseType(e.BaseType())
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

func (tc *TypeChecker) inferFieldExp(
	e *parser.FieldExpContext, line, col int, text string,
) (*tast.FieldExp, error) {
	exp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	fieldName := e.Ident().GetText()

	fieldProviderType, ok := exp.Type().(tast.FieldProvider)
	if !ok {
		return nil, fmt.Errorf(
			"type %s does not have any accessible fields at %d:%d near %s",
			exp.Type(), line, col, text,
		)
	}

	fieldInfo, ok := fieldProviderType.FieldInfo(fieldName)
	if !ok {
		return nil, fmt.Errorf(
			"type %s does not have field %s at %d:%d near %s",
			exp.Type().String(), fieldName, line, col, text,
		)
	}

	return tast.NewFieldExp(
		exp, fieldName, fieldInfo.Type,
		line, col, text,
	), nil
}
