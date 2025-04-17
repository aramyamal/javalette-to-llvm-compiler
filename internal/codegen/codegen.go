package codegen

import (
	"fmt"
	"io"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

type CodeGenerator struct {
	env   *env.Environment[llvm.Reg]
	write *llvm.LLVMWriter
	ng    *NameGenerator
}

func NewCodeGenerator(w io.Writer) *CodeGenerator {
	env := env.NewEnvironment[llvm.Reg]()
	writer := llvm.NewLLVMWriter(w)
	nameGen := NewNameGenerator()
	return &CodeGenerator{env: env, write: writer, ng: nameGen}
}

func (cg *CodeGenerator) GenerateCode(prgm *tast.Prgm) error {
	// boilerplate std functions
	if err := cg.write.Declare(
		llvm.Void, "printInt", llvm.I32); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvm.Void, "printDouble", llvm.Double,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvm.Void, "printString", llvm.I8Ptr,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(llvm.I32, "readInt"); err != nil {
		return err
	}
	if err := cg.write.Declare(llvm.Double, "readDouble"); err != nil {
		return err
	}
	cg.env.EnterContext()
	defer cg.env.ExitContext()

	for _, def := range prgm.Defs {
		cg.ng.resetReg()
		cg.ng.resetLab()

		if err := cg.write.Newline(); err != nil {
			return err
		}
		if err := cg.compileDef(def); err != nil {
			return err
		}
		if err := cg.handleStrings(); err != nil {
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

func (cg *CodeGenerator) emitVarAlloc(
	name string,
	typ llvm.Type,
	init ...llvm.Value,
) error {
	varPtr := llvm.Reg("." + name + "_ptr")
	cg.env.ExtendVar(name, varPtr)
	if err := cg.write.Alloca(varPtr, typ); err != nil {
		return err
	}
	if len(init) > 0 && init[0] != nil {
		if err := cg.write.Store(typ, init[0], varPtr); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CodeGenerator) compileStm(stm tast.Stm) error {
	switch s := stm.(type) {
	case *tast.ExpStm:
		if _, err := cg.compileExp(s.Exp); err != nil {
			return err
		}
		return nil

	case *tast.DeclsStm:
		for _, item := range s.Items {
			llvmType := toLlvmType(item.Type())
			switch i := item.(type) {
			case *tast.NoInitItem:
				if err := cg.emitVarAlloc(i.Id, llvmType); err != nil {
					return err
				}
			case *tast.InitItem:
				value, err := cg.compileExp(i.Exp)
				if err != nil {
					return err
				}
				if err := cg.emitVarAlloc(i.Id, llvmType, value); err != nil {
					return err
				}
			}
		}
		return nil

	case *tast.ReturnStm:
		reg, err := cg.compileExp(s.Exp)
		if err != nil {
			return err
		}
		if err := cg.write.Ret(toLlvmType(s.Type), reg); err != nil {
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

func (cg *CodeGenerator) compileExp(exp tast.Exp) (llvm.Value, error) {
	switch e := exp.(type) {
	case *tast.ParenExp:
		return cg.compileExp(e.Exp)

	case *tast.BoolExp:
		return llvm.LitBool(e.Value), nil

	case *tast.IntExp:
		return llvm.LitInt(e.Value), nil

	case *tast.DoubleExp:
		return llvm.LitDouble(e.Value), nil

	case *tast.StringExp:
		des := cg.ng.nextReg()
		glbVar, strLen := cg.ng.addString(e.Value)
		return des, cg.write.GetElementPtr(
			des,
			llvm.Array(llvm.I8, strLen),
			glbVar,
			0, 0,
		)

	case *tast.IdentExp:
		des := cg.ng.nextReg()
		reg, ok := cg.env.LookupVar(e.Id)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error: undefined variable '%s' encountered"+
					"during code generation at %d:%d near '%s'. "+
					"This should have been caught during type checking.",
				e.Id, e.Line(), e.Col(), e.Text(),
			)
		}
		return des, cg.write.Load(des, toLlvmType(e.Type()), reg)

	case *tast.FuncExp:
		var args []llvm.FuncArg
		for _, exp := range e.Exps {
			value, err := cg.compileExp(exp)
			if err != nil {
				return nil, err
			}
			args = append(args, llvm.Arg(toLlvmType(exp.Type()), value))
		}
		des := cg.ng.nextReg()
		return des, cg.write.Call(
			des,
			toLlvmType(e.Type()),
			llvm.Global(e.Id),
			args...,
		)

	case *tast.NegExp:
		value, err := cg.compileExp(e.Exp)
		if err != nil {
			return nil, err
		}
		des := cg.ng.nextReg()
		llvmType := toLlvmType(e.Type())
		switch llvmType {
		case llvm.I32:
			return des, cg.write.Sub(des, llvmType, llvm.LitInt(0), value)
		case llvm.Double:
			return des, cg.write.Sub(des, llvmType, llvm.LitDouble(0.0), value)
		default:
			return nil, fmt.Errorf(
				"internal compiler error: unable to negate expression "+
					"during code generation at %d:%d near '%s'. "+
					"This should have been caught during type checking.",
				e.Line(), e.Col(), e.Text(),
			)
		}
	case *tast.NotExp:
		value, err := cg.compileExp(e.Exp)
		if err != nil {
			return nil, err
		}
		des := cg.ng.nextReg()
		return des, cg.write.Xor(des, llvm.I1, value, llvm.LitBool(true))

	default:
		return nil, fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) handleStrings() error {
	for name, str := range cg.ng.strMap {
		typ := llvm.Array(llvm.I8, len(str)+1)
		if err := cg.write.Newline(); err != nil {
			return err
		}
		if err := cg.write.InternalConstant(name, typ, str); err != nil {
			return err
		}
	}
	cg.ng.resetStrings()
	return nil
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
