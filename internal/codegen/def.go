package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

func (cg *CodeGenerator) compileDef(def tast.Def) error {
	cg.env.EnterContext()
	switch d := def.(type) {
	case *tast.FuncDef:
		params, err := extractParams(d.Args)
		if err != nil {
			return err
		}
		if err := cg.write.StartDefine(
			toLlvmType(d.Type()),
			llvm.Global(d.Id),
			params...,
		); err != nil {
			return err
		}
		if err := cg.write.Label("entry"); err != nil {
			return err
		}
		for _, param := range params {
			if err := cg.emitVarAlloc(
				string(param.Name),
				param.Type,
				param.Name,
			); err != nil {
				return err
			}
		}
		for _, stm := range d.Stms {
			if err := cg.compileStm(stm); err != nil {
				return err
			}
		}
		return cg.write.EndDefine()
	default:
		return fmt.Errorf(
			"compileDef: unhandled def type %T at %d:%d near '%s'",
			d, d.Line(), d.Col(), d.Text(),
		)
	}
}

func extractParams(args []tast.Arg) ([]llvm.FuncParam, error) {
	var params []llvm.FuncParam
	for _, arg := range args {
		switch a := arg.(type) {
		case *tast.ParamArg:
			params = append(params, llvm.Param(toLlvmType(a.Type()), a.Id))
		default:
			return nil, fmt.Errorf(
				"extractParams: unhandled Arg type %T at %d:%d near '%s'",
				arg, arg.Line(), arg.Col(), arg.Text(),
			)

		}
	}
	return params, nil
}
