package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
)

type TypeChecker struct {
	env *env.Environment[tast.Type]
}

func NewTypeChecker() *TypeChecker {
	env := env.NewEnvironment[tast.Type]()
	return &TypeChecker{env: env}
}

func (tc *TypeChecker) Typecheck(tree parser.IPrgmContext) (*tast.Prgm, error) {
	prgm, ok := tree.(*parser.PrgmContext)
	if !ok {
		return nil, fmt.Errorf("expected *parser.ProgramContext, got %T", tree)
	}
	defs := prgm.AllDef()
	if err := validateMainFunc(defs); err != nil {
		return nil, err
	}

	tc.env.AddStdFunc("printInt", tast.Void, tast.Int)
	tc.env.AddStdFunc("printDouble", tast.Void, tast.Double)
	tc.env.AddStdFunc("printString", tast.Void, tast.String)
	tc.env.AddStdFuncNoParam("readInt", tast.Int)
	tc.env.AddStdFuncNoParam("readDouble", tast.Double)

	tc.env.EnterContext()

	if err := validateFuncSigns(tc.env, defs); err != nil {
		return nil, err
	}

	typedDefs, err := tc.checkDefs(prgm.AllDef())
	if err != nil {
		return nil, err
	}

	tc.env.ExitContext()

	typedPrgm := tast.NewPrgm(typedDefs)
	return typedPrgm, nil
}
