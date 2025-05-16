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
		cg.write.InternalConstant(glbVar, typ, llvmgen.LitString(e.Value))
	}
	cg.write.GetElementPtr(des, llvmgen.Array(llvmgen.I8, strLen),
		llvmgen.Array(llvmgen.I8, strLen).Ptr(), glbVar,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	)
	return des, nil
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
	typ := toLlvmType(e.Type())

	if _, isStruct := typ.(*llvmgen.StructType); isStruct {
		// for arrays and structs, load a pointer to them
		ptrType := typ.Ptr()
		cg.write.Load(des, ptrType, ptrType.Ptr(), reg)
		return des, nil
	} else {
		// for primitive types, load the value
		cg.write.Load(des, typ, typ.Ptr(), reg)
		return des, nil
	}
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
		args = append(args, llvmgen.Arg(toLlvmRetType(exp.Type()), value))
	}

	des := cg.ng.nextReg()
	cg.write.Call(des, toLlvmRetType(e.Type()), llvmgen.Global(e.Id), args...)
	return des, nil
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
		cg.write.Sub(des, llvmType, llvmgen.LitInt(0), value)
	case llvmgen.Double:
		cg.write.Sub(des, llvmType, llvmgen.LitDouble(0.0), value)
	default:
		return nil, fmt.Errorf(
			"internal compiler error: unable to negate expression "+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
}

func (cg *CodeGenerator) compileNotExp(e *tast.NotExp) (llvmgen.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	cg.write.Xor(des, llvmgen.I1, value, llvmgen.LitBool(true))
	return des, nil
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
	cg.write.Load(orig, typ, typ.Ptr(), ptrName)
	incrm := cg.ng.nextReg()

	switch e.Op {
	case tast.OpInc:
		cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1))
	case tast.OpDec:
		cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1))
	default:
		return nil, fmt.Errorf(
			"compileExp->PostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	cg.write.Store(typ, incrm, typ.Ptr(), ptrName)
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
	cg.write.Load(orig, typ, typ.Ptr(), ptrName)
	incrm := cg.ng.nextReg()

	switch e.Op {
	case tast.OpInc:
		cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1))
	case tast.OpDec:
		cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1))
	default:
		return nil, fmt.Errorf(
			"compileExp->PostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	cg.write.Store(typ, incrm, typ.Ptr(), ptrName)
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
		cg.write.Mul(des, toLlvmType(e.Type()), lhs, rhs)
	case tast.OpDiv:
		cg.write.Div(des, toLlvmType(e.Type()), lhs, rhs)
	case tast.OpMod:
		cg.write.Rem(des, toLlvmType(e.Type()), lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->MulExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
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
		cg.write.Add(des, toLlvmType(e.Type()), lhs, rhs)
	case tast.OpSub:
		cg.write.Sub(des, toLlvmType(e.Type()), lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->AddExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
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
		cg.write.CmpLt(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	case tast.OpGt:
		cg.write.CmpGt(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	case tast.OpLe:
		cg.write.CmpLe(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	case tast.OpGe:
		cg.write.CmpGe(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	case tast.OpEq:
		cg.write.CmpEq(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	case tast.OpNe:
		cg.write.CmpNe(des, toLlvmType(e.LeftExp.Type()), lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->CmpExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
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

	cg.write.BrIf(llvmType, lhs, evalLab, falseLab)

	cg.write.Label(evalLab)
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	cg.write.Br(endLab)

	cg.write.Label(falseLab)
	cg.write.Br(endLab)

	des := cg.ng.nextReg()
	cg.write.Label(endLab)
	cg.write.Phi(
		des, llvmType,
		llvmgen.Phi(llvmgen.LitBool(false), falseLab),
		llvmgen.Phi(rhs, evalLab),
	)
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

	cg.write.BrIf(llvmType, lhs, trueLab, evalLab)

	cg.write.Label(evalLab)
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	cg.write.Br(endLab)

	cg.write.Label(trueLab)
	cg.write.Br(endLab)

	des := cg.ng.nextReg()
	cg.write.Label(endLab)
	cg.write.Phi(
		des, llvmType,
		llvmgen.Phi(llvmgen.LitBool(true), trueLab), llvmgen.Phi(rhs, evalLab),
	)
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
	typ := toLlvmRetType(e.Type())
	cg.write.Store(typ, value, typ.Ptr(), ptr)
	return value, nil
}
