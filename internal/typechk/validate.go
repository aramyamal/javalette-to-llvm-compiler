package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
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

	if len(mainFunc.AllArg()) != 0 {
		return fmt.Errorf("entrypoint 'main' may not have input variables")
	}

	if typ, err := toIrType(mainFunc.Type_()); err != nil {
		return err
	} else if typ != types.Int {
		return fmt.Errorf("'main' entrypoint function does not have type int")
	}

	return nil
}

func validateFuncSigns(
	env *env.Environment[types.Type],
	defs []parser.IDefContext,
) error {
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			name := funcDef.Ident().GetText()
			returnType, err := toIrType(funcDef.Type_())
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
