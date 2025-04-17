package typechk

import (
	"fmt"
	"slices"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

func (tc *TypeChecker) checkDefs(defs []parser.IDefContext) ([]tast.Def, error) {

	var typedDefs []tast.Def
	for _, def := range defs {
		typedDef, err := tc.checkDef(def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
		tc.env.SetReturnType(types.Unknown)
	}
	return typedDefs, nil
}

func (tc *TypeChecker) checkDef(def parser.IDefContext) (tast.Def, error) {
	tc.env.EnterContext()
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		return tc.checkFuncDef(d, line, col, text)
	default:
		return nil, fmt.Errorf(
			"checkDef: unhandled def type %T at %d:%d near '%s'",
			d, line, col, text,
		)
	}
}

func (tc *TypeChecker) checkFuncDef(
	d *parser.FuncDefContext, line, col int, text string,
) (*tast.FuncDef, error) {
	_, params, err := extractParams(d.AllArg())
	if err != nil {
		return nil, err
	}

	for varName, typ := range params {
		ok := tc.env.ExtendVar(varName, typ)
		if !ok {
			return nil, fmt.Errorf(
				"duplicate parameter name '%s' in function '%s' at %d:%d",
				varName, d.Ident().GetText(), line, col,
			)
		}
	}

	typ, err := toIrType(d.Type_())
	if err != nil {
		return nil, err
	}
	tc.env.SetReturnType(typ)

	var typedStms []tast.Stm
	for _, stm := range d.AllStm() {
		typedStm, err := tc.checkStm(stm)
		if err != nil {
			return nil, err
		}
		typedStms = append(typedStms, typedStm)
	}

	hasReturn := slices.ContainsFunc(typedStms, guaranteesReturn)

	if typ != types.Void && !hasReturn {
		return nil, fmt.Errorf(
			"function '%s' at %d:%d does not have a return", text, line, col,
		)
	}

	typedArgs, err := toAstArgs(d.AllArg())
	if err != nil {
		return nil, err
	}
	tc.env.ExitContext()
	return tast.NewFuncDef(
		d.Ident().GetText(),
		typedArgs,
		typedStms,
		typ,
		line, col, text,
	), nil

}
