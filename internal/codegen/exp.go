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

	case *tast.PostExp:

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

	case *tast.PreExp:

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

	default:
		return nil, fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}
