package codegen

import (
	"fmt"
	"io"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

type CodeGenerator struct {
	env   *env.Environment[llvm.Var]
	write *llvm.LLVMWriter
	ng    *NameGenerator
}

func NewCodeGenerator(w io.Writer) *CodeGenerator {
	env := env.NewEnvironment[llvm.Var]()
	writer := llvm.NewLLVMWriter(w)
	nameGen := &NameGenerator{}
	return &CodeGenerator{env: env, write: writer, ng: nameGen}
}

func (cg *CodeGenerator) GenerateCode(prgm *tast.Prgm) error {
	// boilerplate std functions
	if err := cg.write.Declare(types.Void, "printInt", types.Int); err != nil {
		return err
	}
	if err := cg.write.Declare(
		types.Void, "printDouble", types.Double,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(
		types.Void, "printString", types.String,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(types.Int, "readInt"); err != nil {
		return err
	}
	if err := cg.write.Declare(types.Double, "readDouble"); err != nil {
		return err
	}
	if err := cg.write.Newline(); err != nil {
		return err
	}

	cg.env.EnterContext()
	defer cg.env.ExitContext()

	for _, def := range prgm.Defs {
		cg.ng.resetReg()
		cg.ng.resetLab()
		if err := cg.compileDef(def); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CodeGenerator) compileDef(def tast.Def) error {
	cg.env.EnterContext()
	switch d := def.(type) {
	case *tast.FuncDef:
		params, err := extractParams(d.Args)
		if err != nil {
			return err
		}
		if err := cg.write.StartDefine(
			d.Type(),
			llvm.Global(d.Id),
			params...,
		); err != nil {
			return err
		}
		if err := cg.write.Label("entry"); err != nil {
			return err
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

func (cg *CodeGenerator) compileStm(stm tast.Stm) error {
	switch s := stm.(type) {
	case *tast.ExpStm:
		if err := cg.compileExp(s.Exp); err != nil {
			return err
		}
		return nil
	case *tast.ReturnStm:
		if err := cg.compileExp(s.Exp); err != nil {
			return err
		}
		if err := cg.write.Ret(s.Type, cg.ng.currentReg()); err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf(
			"compileStm: unhandled stm type %T at %d:%d near '%s'",
			s, s.Line(), s.Col(), s.Text(),
		)
	}
}

func (cg *CodeGenerator) compileExp(exp tast.Exp) error {
	switch e := exp.(type) {
	case *tast.ParenExp:
		return cg.compileExp(e.Exp)
	case *tast.BoolExp:
		return cg.write.Constant(
			cg.ng.nextReg(),
			types.Bool,
			llvm.LitBool(e.Value),
		)
	case *tast.IntExp:
		return cg.write.Constant(
			cg.ng.nextReg(),
			types.Int,
			llvm.LitInt(e.Value),
		)
	case *tast.DoubleExp:
		return cg.write.Constant(
			cg.ng.nextReg(),
			types.Double,
			llvm.LitDouble(e.Value),
		)
	case *tast.IdentExp:
		// varName, ok := cg.env.LookupVar(e.Id)
		return nil
	default:
		return fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}

func extractParams(args []tast.Arg) ([]llvm.Param, error) {
	var params []llvm.Param
	for _, arg := range args {
		switch a := arg.(type) {
		case *tast.ParamArg:
			params = append(params, llvm.NewParam(a.Type(), a.Id))
		default:
			return nil, fmt.Errorf(
				"extractParams: unhandled Arg type %T at %d:%d near '%s'",
				arg, arg.Line(), arg.Col(), arg.Text(),
			)

		}
	}
	return params, nil
}
