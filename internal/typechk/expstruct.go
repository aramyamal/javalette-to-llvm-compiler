package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) inferNullPtrExp(
	e *parser.NullPtrExpContext, line, col int, text string,
) (*tast.NullPtrExp, error) {
	typ, err := tc.toTastType(e.Type_())
	if err != nil {
		return nil, err
	}
	return tast.NewNullPtrExp(typ, line, col, text), nil
}

func (tc *TypeChecker) inferNewStructExp(
	e *parser.NewStructExpContext, line, col int, text string,
) (*tast.NewStructExp, error) {
	name := e.Ident().GetText()
	// try typedef first
	if typ, ok := tc.env.LookupTypedef(name); ok {
		return tast.NewNewStructExp(typ, line, col, text), nil
	}
	// fallback: try struct
	if typ, ok := tc.env.LookupStruct(name); ok {
		return tast.NewNewStructExp(tast.Pointer(typ), line, col, text), nil
	}
	return nil, fmt.Errorf(
		"cannot allocate new struct '%s' if it has not been defined, at "+
			"%d:%d near %s", name, line, col, text,
	)
}

func (tc *TypeChecker) inferDerefExp(
	e *parser.DerefExpContext, line, col int, text string,
) (*tast.DerefExp, error) {
	exp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}

	fieldName := e.Ident().GetText()

	// unwrap typedefs before pointer check
	typ := UnwrapTypedef(exp.Type())

	pointerType, ok := typ.(*tast.PointerType)
	if !ok {
		return nil, fmt.Errorf(
			"type %s is not a pointer to be able to have fields that can be "+
				"dereferenced at %d:%d near %s", exp.Type(), line, col, text,
		)
	}

	// unwrap typedefs before struct/field check
	elemType := UnwrapTypedef(pointerType.Elem)

	fieldProviderType, ok := elemType.(tast.FieldProvider)
	if !ok {
		return nil, fmt.Errorf(
			"type that type %s points to does not have any accessible fields"+
				" at %d:%d near %s", exp.Type(), line, col, text,
		)
	}

	fieldInfo, ok := fieldProviderType.FieldInfo(fieldName)
	if !ok {
		return nil, fmt.Errorf(
			"type %s does not have field %s at %d:%d near %s",
			exp.Type().String(), fieldName, line, col, text,
		)
	}

	return tast.NewDerefExp(
		exp, fieldName, fieldInfo.Type,
		line, col, text,
	), nil
}

func UnwrapTypedef(t tast.Type) tast.Type {
	for {
		if td, ok := t.(*tast.TypedefType); ok {
			t = td.Aliased
		} else {
			return t
		}
	}
}
