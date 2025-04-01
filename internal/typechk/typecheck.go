package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
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

	env := NewEnvironment[tast.Type]()

	env.AddStdFunc("printInt", tast.Void, tast.Int)
	env.AddStdFunc("printDouble", tast.Void, tast.Double)
	env.AddStdFunc("printString", tast.Void, tast.String)
	env.AddStdFunc("readInt", tast.Int, tast.Unknown)
	env.AddStdFunc("readDouble", tast.Double, tast.Unknown)

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
