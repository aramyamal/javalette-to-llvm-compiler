package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileStringExp(e *tast.StringExp) (
	llvmgen.Value, error,
) {
	des := cg.ng.nextReg()
	glbVar, strLen, alreadyWritten := cg.ng.getOrAddString(e.Value)
	if !alreadyWritten {
		typ := llvmgen.Array(llvmgen.I8, strLen)
		if err := cg.write.InternalConstant(
			glbVar, typ, llvmgen.LitString(e.Value),
		); err != nil {
			return nil, err
		}
	}
	return des, cg.write.GetElementPtr(
		des,
		llvmgen.Array(llvmgen.I8, strLen),
		glbVar,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	)
}

func (cg *CodeGenerator) compileIdentExp(e *tast.IdentExp) (
	llvmgen.Value, error,
) {
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

func (cg *CodeGenerator) compileFuncExp(e *tast.FuncExp) (
	llvmgen.Value, error,
) {
	var args []llvmgen.FuncArg
	for _, exp := range e.Exps {
		value, err := cg.compileExp(exp)
		if err != nil {
			return nil, err
		}
		args = append(args, llvmgen.Arg(toLlvmType(exp.Type()), value))
	}
	des := cg.ng.nextReg()
	return des, cg.write.Call(
		des,
		toLlvmType(e.Type()),
		llvmgen.Global(e.Id),
		args...,
	)
}

func (cg *CodeGenerator) compileFieldExp(e *tast.FieldExp) (
	llvmgen.Value, error,
) {
	return nil, fmt.Errorf("compileFieldExp: not yet implemented")
}

func (cg *CodeGenerator) compileNegExp(e *tast.NegExp) (llvmgen.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	llvmType := toLlvmType(e.Type())
	switch llvmType {
	case llvmgen.I32:
		return des, cg.write.Sub(des, llvmType, llvmgen.LitInt(0), value)
	case llvmgen.Double:
		return des, cg.write.Sub(des, llvmType, llvmgen.LitDouble(0.0), value)
	default:
		return nil, fmt.Errorf(
			"internal compiler error: unable to negate expression "+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileNotExp(e *tast.NotExp) (llvmgen.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	return des, cg.write.Xor(des, llvmgen.I1, value, llvmgen.LitBool(true))
}

func (cg *CodeGenerator) compilePostExp(e *tast.PostExp) (
	llvmgen.Value, error,
) {
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
	case tast.OpInc:
		if err := cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1)); err != nil {
			return nil, err
		}
	case tast.OpDec:
		if err := cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1)); err != nil {
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

func (cg *CodeGenerator) compilePreExp(e *tast.PreExp) (llvmgen.Value, error) {
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
	case tast.OpInc:
		if err := cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1)); err != nil {
			return nil, err
		}
	case tast.OpDec:
		if err := cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1)); err != nil {
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

func (cg *CodeGenerator) compileMulExp(e *tast.MulExp) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	switch e.Op {
	case tast.OpMul:
		if err := cg.write.Mul(des, toLlvmType(e.Type()), lhs, rhs); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpDiv:
		if err := cg.write.Div(des, toLlvmType(e.Type()), lhs, rhs); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpMod:
		if err := cg.write.Rem(des, toLlvmType(e.Type()), lhs, rhs); err != nil {
			return nil, err
		}
		return des, nil
	default:
		return nil, fmt.Errorf(
			"compileExp->MulExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileAddExp(e *tast.AddExp) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	switch e.Op {
	case tast.OpAdd:
		if err := cg.write.Add(des, toLlvmType(e.Type()), lhs, rhs); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpSub:
		if err := cg.write.Sub(des, toLlvmType(e.Type()), lhs, rhs); err != nil {
			return nil, err
		}
		return des, nil
	default:
		return nil, fmt.Errorf(
			"compileExp->AddExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileCmpExp(e *tast.CmpExp) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	switch e.Op {
	case tast.OpLt:
		if err := cg.write.CmpLt(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpGt:
		if err := cg.write.CmpGt(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpLe:
		if err := cg.write.CmpLe(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpGe:
		if err := cg.write.CmpGe(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpEq:
		if err := cg.write.CmpEq(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	case tast.OpNe:
		if err := cg.write.CmpNe(
			des, toLlvmType(e.LeftExp.Type()), lhs, rhs,
		); err != nil {
			return nil, err
		}
		return des, nil
	default:
		return nil, fmt.Errorf(
			"compileExp->CmpExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileAndExp(e *tast.AndExp) (llvmgen.Value, error) {
	llvmType := toLlvmType(e.Type())
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	falseLab := cg.ng.nextLab()
	evalLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()

	if err := cg.write.BrIf(llvmType, lhs, evalLab, falseLab); err != nil {
		return nil, err
	}

	if err := cg.write.Label(evalLab); err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	if err := cg.write.Br(endLab); err != nil {
		return nil, err
	}

	if err := cg.write.Label(falseLab); err != nil {
		return nil, err
	}
	if err := cg.write.Br(endLab); err != nil {
		return nil, err
	}

	des := cg.ng.nextReg()
	if err := cg.write.Label(endLab); err != nil {
		return nil, err
	}
	if err := cg.write.Phi(
		des,
		llvmType,
		llvmgen.Phi(llvmgen.LitBool(false), falseLab),
		llvmgen.Phi(rhs, evalLab),
	); err != nil {
		return nil, err
	}
	return des, nil
}

func (cg *CodeGenerator) compileOrExp(e *tast.OrExp) (llvmgen.Value, error) {
	llvmType := toLlvmType(e.Type())
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	trueLab := cg.ng.nextLab()
	evalLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()

	if err := cg.write.BrIf(llvmType, lhs, trueLab, evalLab); err != nil {
		return nil, err
	}

	if err := cg.write.Label(evalLab); err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	if err := cg.write.Br(endLab); err != nil {
		return nil, err
	}

	if err := cg.write.Label(trueLab); err != nil {
		return nil, err
	}
	if err := cg.write.Br(endLab); err != nil {
		return nil, err
	}

	des := cg.ng.nextReg()
	if err := cg.write.Label(endLab); err != nil {
		return nil, err
	}
	if err := cg.write.Phi(
		des,
		llvmType,
		llvmgen.Phi(llvmgen.LitBool(true), trueLab),
		llvmgen.Phi(rhs, evalLab),
	); err != nil {
		return nil, err
	}
	return des, nil
}

func (cg *CodeGenerator) compileAssignExp(
	e *tast.AssignExp,
) (llvmgen.Value, error) {
	ptr, ok := cg.env.LookupVar(e.Id)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error: undefined variable '%s' encountered"+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Id, e.Line(), e.Col(), e.Text(),
		)
	}
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	if err := cg.write.Store(toLlvmType(e.Type()), value, ptr); err != nil {
		return nil, err
	}
	return value, nil
}
