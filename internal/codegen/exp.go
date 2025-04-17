package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

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
		return cg.compileStringExp(e)
	case *tast.IdentExp:
		return cg.compileIdentExp(e)
	case *tast.FuncExp:
		return cg.compileFuncExp(e)
	case *tast.NegExp:
		return cg.compileNegExp(e)
	case *tast.NotExp:
		return cg.compileNotExp(e)
	case *tast.PostExp:
		return cg.compilePostExp(e)
	case *tast.PreExp:
		return cg.compilePreExp(e)
	default:
		return nil, fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileStringExp(
	e *tast.StringExp,
) (llvm.Value, error) {
	des := cg.ng.nextReg()
	glbVar, strLen := cg.ng.addString(e.Value)
	return des, cg.write.GetElementPtr(
		des,
		llvm.Array(llvm.I8, strLen),
		glbVar,
		0, 0,
	)
}

func (cg *CodeGenerator) compileIdentExp(
	e *tast.IdentExp,
) (llvm.Value, error) {
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
}

func (cg *CodeGenerator) compileFuncExp(
	e *tast.FuncExp,
) (llvm.Value, error) {
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
}

func (cg *CodeGenerator) compileNegExp(
	e *tast.NegExp,
) (llvm.Value, error) {
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
}

func (cg *CodeGenerator) compileNotExp(
	e *tast.NotExp,
) (llvm.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	return des, cg.write.Xor(des, llvm.I1, value, llvm.LitBool(true))
}

func (cg *CodeGenerator) compilePostExp(
	e *tast.PostExp,
) (llvm.Value, error) {
	ptrName, ok := cg.env.LookupVar(e.Id)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error: undefined variable '%s' encountered"+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Id, e.Line(), e.Col(), e.Text(),
		)
	}
	orig := cg.ng.nextReg()
	typ := toLlvmType(e.Type())
	if err := cg.write.Load(orig, typ, ptrName); err != nil {
		return nil, err
	}
	incrm := cg.ng.nextReg()

	switch e.Op {
	case types.OpInc:
		if err := cg.write.Add(incrm, typ, llvm.LitInt(1), orig); err != nil {
			return nil, err
		}
	case types.OpDec:
		if err := cg.write.Sub(incrm, typ, llvm.LitInt(1), orig); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(
			"compileExp->PostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	if err := cg.write.Store(typ, incrm, ptrName); err != nil {
		return nil, err
	}
	return orig, nil
}

func (cg *CodeGenerator) compilePreExp(
	e *tast.PreExp,
) (llvm.Value, error) {
	ptrName, ok := cg.env.LookupVar(e.Id)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error: undefined variable '%s' encountered"+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Id, e.Line(), e.Col(), e.Text(),
		)
	}
	orig := cg.ng.nextReg()
	typ := toLlvmType(e.Type())
	if err := cg.write.Load(orig, typ, ptrName); err != nil {
		return nil, err
	}
	incrm := cg.ng.nextReg()

	switch e.Op {
	case types.OpInc:
		if err := cg.write.Add(incrm, typ, llvm.LitInt(1), orig); err != nil {
			return nil, err
		}
	case types.OpDec:
		if err := cg.write.Sub(incrm, typ, llvm.LitInt(1), orig); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(
			"compileExp->PostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	if err := cg.write.Store(typ, incrm, ptrName); err != nil {
		return nil, err
	}
	return incrm, nil
}
