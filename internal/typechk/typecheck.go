package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

type TypeChecker struct {
	env *env.Environment[types.Type]
}

func NewTypeChecker() *TypeChecker {
	env := env.NewEnvironment[types.Type]()
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

	tc.env.AddStdFunc("printInt", types.Void, types.Int)
	tc.env.AddStdFunc("printDouble", types.Void, types.Double)
	tc.env.AddStdFunc("printString", types.Void, types.String)
	tc.env.AddStdFuncNoParam("readInt", types.Int)
	tc.env.AddStdFuncNoParam("readDouble", types.Double)

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
