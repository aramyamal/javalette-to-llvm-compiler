package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func checkDefs(
	env *Environment[tast.Type],
	defs []parser.IDefContext,
) ([]tast.Def, error) {

	var typedDefs []tast.Def
	for _, def := range defs {
		typedDef, err := checkDef(env, def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
	}
	return typedDefs, nil
}

func checkDef(
	env *Environment[tast.Type],
	def parser.IDefContext,
) (tast.Def, error) {
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		// TODO:
		// handle Ident by adding to func. context,
		// handle args by adding to environment,

		hasReturn := false
		var typedStms []tast.Stm
		for _, stm := range d.AllStm() {
			typedStm, err := checkStm(env, stm)
			if err != nil {
				return nil, err
			}
			if _, ok := typedStm.(*tast.ReturnStm); ok {
				hasReturn = true
			}
			typedStms = append(typedStms, typedStm)
		}

		typ, err := toAstType(d.Type_())
		if err != nil {
			return nil, err
		}

		if typ != tast.Void && !hasReturn {
			return nil, fmt.Errorf(
				"function '%s' at %d:%d does not have a return",
				text, line, col,
			)
		}

		typedArgs, err := toAstArgs(d.AllArg())
		if err != nil {
			return nil, err
		}
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
