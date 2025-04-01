package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func validateMainFunc(defs []parser.IDefContext) error {
	var mainFunc *parser.FuncDefContext
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			if funcDef.Ident().GetText() == "main" {
				mainFunc = funcDef
				break
			}
		}
	}

	if mainFunc == nil {
		return fmt.Errorf("program has no entrypoint 'main'")
	}

	if typ, err := toAstType(mainFunc.Type_()); err != nil {
		return err
	} else if typ != tast.Int {
		return fmt.Errorf("'main' entrypoint function does not have type int")
	}

	// TODO add check that main does not have params

	return nil
}

func validateFuncSigns(
	env *Environment[tast.Type],
	defs []parser.IDefContext,
) error {
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			name := funcDef.Ident().GetText()
			returnType, err := toAstType(funcDef.Type_())
			if err != nil {
				return err
			}

			paramNames, params, err := extractParams(funcDef.AllArg())
			if err != nil {
				return err
			}

			if ok := env.ExtendFunc(name, paramNames, params, returnType); !ok {
				return fmt.Errorf(
					"redefinition of function '%s' at %d:%d",
					name, def.GetStart().GetLine(), def.GetStart().GetColumn(),
				)
			}
		}
	}
	return nil
}
