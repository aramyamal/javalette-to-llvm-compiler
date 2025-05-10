package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) checkForEachStm(
	s *parser.ForEachStmContext, line, col int, text string,
) (*tast.ForEachStm, error) {

	tc.env.EnterContext()
	defer tc.env.ExitContext()

	typ, err := toTastType(s.Type_())
	if err != nil {
		return nil, err
	}

	id := s.Ident().GetText()
	if ok := tc.env.ExtendVar(id, typ); !ok {
		return nil, fmt.Errorf(
			"redefinition of 'for each'-variable %s %d:%d near %s",
			id, line, col, text,
		)
	}

	exp, err := tc.inferExp(s.Exp())
	if err != nil {
		return nil, err
	}

	arrType, ok := exp.Type().(*tast.ArrayType)
	if !ok {
		return nil, fmt.Errorf(
			"can only iterate over array objects at %d:%d near %s",
			line, col, text,
		)
	}

	if !isConvertible(typ, arrType.Elem) {
		return nil, fmt.Errorf(
			"for-each variable %s has type %s, but array elements have type %s"+
				" at %d:%d near %s",
			id, typ.String(), arrType.Elem.String(), line, col, text,
		)
	}

	stm, err := tc.checkStm(s.Stm())
	if err != nil {
		return nil, err
	}

	return tast.NewForEachStm(typ, id, exp, stm, line, col, text), nil
}
