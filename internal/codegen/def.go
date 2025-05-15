package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileDef(def tast.Def) error {
	cg.env.EnterContext()
	switch d := def.(type) {
	case *tast.FuncDef:
		return cg.compileFuncDef(d)
	default:
		return fmt.Errorf(
			"compileDef: unhandled def type %T at %d:%d near '%s'",
			d, d.Line(), d.Col(), d.Text(),
		)
	}
}

func (cg *CodeGenerator) compileFuncDef(d *tast.FuncDef) error {
	params, err := extractParams(d.Args)
	if err != nil {
		return err
	}
	cg.write.StartDefine(toLlvmRetType(d.Type()), llvmgen.Global(d.Id), params...)
	cg.write.Label("entry")
	for _, param := range params {
		cg.emitVarAlloc(string(param.Name), param.Type, param.Name)
	}
	for _, stm := range d.Stms {
		if err := cg.compileStm(stm); err != nil {
			return err
		}
	}
	return cg.write.EndDefine()
}

func extractParams(args []tast.Arg) ([]llvmgen.FuncParam, error) {
	var params []llvmgen.FuncParam
	for _, arg := range args {
		switch a := arg.(type) {
		case *tast.ParamArg:
			params = append(params, llvmgen.Param(toLlvmRetType(a.Type()), a.Id))
		default:
			return nil, fmt.Errorf(
				"extractParams: unhandled Arg type %T at %d:%d near '%s'",
				arg, arg.Line(), arg.Col(), arg.Text(),
			)
		}
	}
	return params, nil
}
