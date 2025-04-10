package typechk

import (
	"fmt"
	"slices"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/ir"
)

func checkDefs(
	env *env.Environment[ir.Type],
	defs []parser.IDefContext,
) ([]tast.Def, error) {

	var typedDefs []tast.Def
	for _, def := range defs {
		typedDef, err := checkDef(env, def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
		env.SetReturnType(ir.Unknown)
	}
	return typedDefs, nil
}

func checkDef(
	env *env.Environment[ir.Type],
	def parser.IDefContext,
) (tast.Def, error) {
	env.EnterContext()
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		// TODO:
		// handle Ident by adding to func. context,

		_, params, err := extractParams(d.AllArg())
		if err != nil {
			return nil, err
		}

		for varName, typ := range params {
			ok := env.ExtendVar(varName, typ)
			if !ok {
				return nil, fmt.Errorf(
					"duplicate parameter name '%s' in function '%s' at %d:%d",
					varName, d.Ident().GetText(), line, col,
				)
			}
		}

		typ, err := toAstType(d.Type_())
		if err != nil {
			return nil, err
		}
		env.SetReturnType(typ)

		var typedStms []tast.Stm
		for _, stm := range d.AllStm() {
			typedStm, err := checkStm(env, stm)
			if err != nil {
				return nil, err
			}
			typedStms = append(typedStms, typedStm)
		}

		hasReturn := slices.ContainsFunc(typedStms, guaranteesReturn)

		if typ != ir.Void && !hasReturn {
			return nil, fmt.Errorf(
				"function '%s' at %d:%d does not have a return",
				text, line, col,
			)
		}

		typedArgs, err := toAstArgs(d.AllArg())
		if err != nil {
			return nil, err
		}
		env.ExitContext()
		return tast.NewFuncDef(
			d.Ident().GetText(),
			typedArgs,
			typedStms,
			typ,
			line, col, text,
		), nil
	default:
		return nil, fmt.Errorf(
			"checkDef: unhandled def type %T at %d:%d near '%s'",
			d, line, col, text,
		)
	}
}
