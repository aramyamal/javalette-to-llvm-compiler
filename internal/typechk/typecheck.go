package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/ir"
)

func Typecheck(tree parser.IPrgmContext) (*tast.Prgm, error) {
	prgm, ok := tree.(*parser.PrgmContext)
	if !ok {
		return nil, fmt.Errorf("expected *parser.ProgramContext, got %T", tree)
	}
	defs := prgm.AllDef()
	if err := validateMainFunc(defs); err != nil {
		return nil, err
	}

	env := env.NewEnvironment[ir.Type]()

	env.AddStdFunc("printInt", ir.Void, ir.Int)
	env.AddStdFunc("printDouble", ir.Void, ir.Double)
	env.AddStdFunc("printString", ir.Void, ir.String)
	env.AddStdFuncNoParam("readInt", ir.Int)
	env.AddStdFuncNoParam("readDouble", ir.Double)

	env.EnterContext()

	if err := validateFuncSigns(env, defs); err != nil {
		return nil, err
	}

	typedDefs, err := checkDefs(env, prgm.AllDef())
	if err != nil {
		return nil, err
	}

	env.ExitContext()

	typedPrgm := tast.NewPrgm(typedDefs)
	return typedPrgm, nil
}
