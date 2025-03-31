package typechecker

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func checkDefs(
	env *Environment[typedast.Type],
	defs []parser.IDefContext,
) ([]typedast.Def, error) {

	var typedDefs []typedast.Def
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
	env *Environment[typedast.Type],
	def parser.IDefContext,
) (typedast.Def, error) {
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		// TODO:
		// handle Ident by adding to func. context,
		// handle args by adding to environment,

		var typedStms []typedast.Stm
		for _, stm := range d.AllStm() {
			typedStm, err := checkStm(env, stm)
			if err != nil {
				return nil, err
			}
			typedStms = append(typedStms, typedStm)
		}
		typ, err := toAstType(d.Type_())
		if err != nil {
			return nil, err
		}
		typedArgs, err := toAstArgs(d.AllArg())
		if err != nil {
			return nil, err
		}
		return typedast.NewFuncDef(
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
